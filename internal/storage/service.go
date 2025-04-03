package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fileupbackendv2/config"
	"fileupbackendv2/internal/auth"
	"fileupbackendv2/internal/db"
	"fileupbackendv2/internal/models"
	"fileupbackendv2/internal/storage/image"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader, username, folderPath string) error {
	filenameCleand := strings.Replace(fileHeader.Filename, " ", "_", -1)
	pathFromUploadsDir := filepath.Join(username, folderPath, filenameCleand)
	relativePath := filepath.Join(config.UploadsDir, pathFromUploadsDir)
	filePath := filepath.Join(config.BaseDir, relativePath)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	tx := db.StorageDB.Begin()
	user, err := models.GetUserByUsername(username)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmd := models.FileMetadata{
		Name:     fileHeader.Filename,
		FilePath: pathFromUploadsDir,
		UserID:   user.ID,
		SizeInMb: float64(fileHeader.Size) / 1024 / 1024,
	}
	err = tx.Create(&fmd).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	reader := bufio.NewReader(file)
	writer := io.Writer(f)
	_, err = io.Copy(writer, reader)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetSignedUrl(filePath string, userId string) (string, error) {
	var fmd models.FileMetadata
	db.StorageDB.Where("file_path = ? AND user_id = ?", filePath, userId).First(&fmd)

	expiresAt := time.Now().Add(600 * time.Minute).Unix()
	expriresAtStr := fmt.Sprintf("%d", expiresAt)
	content := fmt.Sprintf("%s%s%s", config.SecretKey, expriresAtStr, filePath)

	signature := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%s?signature=%s&expires_at=%d", filePath, hex.EncodeToString(signature[:]), expiresAt), nil
}

func GetFiles(userId string) ([]models.FileMetadata, error) {
	var files []models.FileMetadata
	err := db.StorageDB.Where("user_id = ?", userId).Find(&files).Error
	return files, err
}

func DownloadFile(filePath, signature, expiresAt string, compressed bool) ([]byte, error) {

	content := fmt.Sprintf("%s%s%s", config.SecretKey, expiresAt, filePath)
	hash := sha256.Sum256([]byte(content))
	if hex.EncodeToString(hash[:]) != signature {
		return nil, errors.New("invalid signature")
	}
	var absPath string
	if compressed {
		absPath = filepath.Join(config.BaseDir, config.UploadsDir, image.GetCompressedPath(filePath))
	}
	f, err := os.ReadFile(absPath)
	if err != nil {
		absPath = filepath.Join(config.BaseDir, config.UploadsDir, filePath)
		f, err = os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

type DirEntry struct {
	Name      string
	IsDir     bool
	Extension string
	SignedUrl string
	Size      float64
	Path      string
}

func GetSizeForDirEntry(file fs.DirEntry) float64 {
	fi, err := file.Info()
	if err != nil {
		return 0
	}
	return float64(fi.Size() / 1024 / 1024)
}

func FindFiles(apiKey string, folder string) ([]DirEntry, float64, error) {
	user, err := auth.GetUserByKey(apiKey)
	if err != nil {
		return nil, 0, err
	}

	pathFromUploadsDir := filepath.Join(user.Username, folder)
	relativePath := filepath.Join(config.UploadsDir, pathFromUploadsDir)
	folderPath := filepath.Join(config.BaseDir, relativePath)
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, 0, err
	}

	var entries []DirEntry
	for _, file := range files {
		signedUrl, err := GetSignedUrl(pathFromUploadsDir+"/"+file.Name(), user.Username)
		if err != nil {
			signedUrl = ""
		}
		entries = append(entries, DirEntry{
			Name:      file.Name(),
			Path:      folder + "/" + file.Name(),
			IsDir:     file.IsDir(),
			Extension: filepath.Ext(file.Name()),
			SignedUrl: config.APIBaseUrl + "files/download/" + signedUrl,
			Size:      GetSizeForDirEntry(file),
		})
	}
	var folderSize float64
	folderData, err := models.GetDirByPathorName(pathFromUploadsDir, folder, user.Username)
	if err == nil {
		folderSize = folderData.SizeInMb
	}
	return entries, folderSize, nil
}

type FolderEntry struct {
	Name string
	Path string
}

func splitPathTillUserDir(path string, username string) string {
	split := strings.Split(path, "/")
	for i := len(split) - 1; i >= 0; i-- {
		if split[i] == username {
			return filepath.Join(split[i+1:]...)
		}
	}
	return path
}

func GetAllFolders(apiKey string) ([]FolderEntry, error) {
	user, err := auth.GetUserByKey(apiKey)
	if err != nil {
		return nil, err
	}
	pathFromUploadsDir := filepath.Join(user.Username)
	relativePath := filepath.Join(config.UploadsDir, pathFromUploadsDir)
	folderPath := filepath.Join(config.BaseDir, relativePath)
	var subDirs []FolderEntry
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != folderPath {
			subDirs = append(subDirs, FolderEntry{
				Name: info.Name(),
				Path: splitPathTillUserDir(path, user.Username),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return subDirs, nil
}
