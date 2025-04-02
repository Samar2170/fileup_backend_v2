package image

import (
	"fileupbackendv2/config"
	"fileupbackendv2/internal/db"
	"image"
	"image/jpeg"
	"image/png"

	storagemodels "fileupbackendv2/internal/storage/storageModels"
	"fileupbackendv2/internal/utils"
	"fileupbackendv2/pkg/logging"
	"log"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
)

const (
	// cwp = compressed webp
	cwp    = "_cwp"
	cwpExt = ".webp"
)

func MarkImages() error {
	var fmds []storagemodels.FileMetadata
	err := db.StorageDB.
		Where("is_image IS NULL OR is_image = ?", false).
		Find(&fmds).Error
	if err != nil {
		return err
	}
	tx := db.StorageDB.Begin()
	for _, fmd := range fmds {
		ext := filepath.Ext(fmd.FilePath)
		if utils.IfArrayContains([]string{".jpg", ".jpeg", ".png"}, ext) {
			fmd.IsImage = true
			err := tx.Save(&fmd).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func GetCompressedPath(finalPath string) string {
	ext := filepath.Ext(finalPath)
	cleanFpath := finalPath[:len(finalPath)-len(ext)]
	cwpFpath := cleanFpath + cwp + cwpExt
	return cwpFpath
}

func CompressImage(fpath string, quality int) error {
	ext := filepath.Ext(fpath)
	var finalPath string
	if config.UploadsDir == fpath[:len(config.UploadsDir)] {
		finalPath = fpath
	} else {
		finalPath = filepath.Join(config.UploadsDir, fpath)
	}
	file, err := os.Open(finalPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cwpFpath := GetCompressedPath(finalPath)
	if _, err := os.Stat(cwpFpath); !os.IsNotExist(err) {
		return nil
	}

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Println("Error decoding image", err)
			return err
		}
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			log.Println("Error decoding image", err)
			return err
		}
	}
	if err != nil {
		return err
	}
	out, err := os.Create(cwpFpath)
	if err != nil {
		return err
	}
	defer out.Close()
	err = webp.Encode(out, img, &webp.Options{
		Quality: float32(quality),
	})
	if err != nil {
		return err
	}
	return nil
}

func CompressImages(quality int) error {
	var fmds []storagemodels.FileMetadata
	err := db.StorageDB.
		Where("compressed_version_available = ? OR compressed_version_available IS NULL", false).
		Where("is_image = ?", true).
		Find(&fmds).Error
	if err != nil {
		return err
	}
	for _, fmd := range fmds {
		err := CompressImage(fmd.FilePath, quality)
		if err != nil {
			logging.Errorlogger.Error().Msgf("Error compressing image: %s", err.Error())
			continue
		}
		fmd.CompressedVersionAvailable = true
		db.StorageDB.Save(&fmd)
	}
	return nil
}

func CompressImagesToWebPForDir(dirPath string, quality int) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg" || filepath.Ext(path) == ".png" {
			return CompressImage(path, quality)
		}
		return nil
	})
	return err
}
