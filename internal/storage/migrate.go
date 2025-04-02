package storage

import (
	"fileupbackendv2/config"
	"fileupbackendv2/internal/db"
	storagemodels "fileupbackendv2/internal/storage/storageModels"
	"fileupbackendv2/pkg/logging"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserAPIKeyED struct {
	*gorm.Model
	Username string
	APIKey   string
	PIN      string `gorm:"column:pin;varchar(6)"`
}

func (UserAPIKeyED) TableName() string {
	return "user_api_keys"
}

type FileMetadataED struct {
	Name     string `gorm:"index"`
	FilePath string `gorm:"index"`
	SizeInMb float64
	IsPublic bool

	IsImage                    bool
	CompressedVersionAvailable bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (FileMetadataED) TableName() string {
	return "file_metadata"
}

//	approach1 -> load from db to db
//
// approach2 -> just migrate the uploads folder
func CheckMigrationSuccessFull() {
	files, err := os.ReadDir(config.BaseDir)
	if err != nil {
		panic(err)
	}
	dbFiles := []string{}
	for _, file := range files {
		extension := filepath.Ext(file.Name())
		if extension == ".db" || extension == ".sqlite" {
			if file.Name() == config.StorageDbFile {
				continue
			}
			dbFiles = append(dbFiles, file.Name())
		}
	}

	for _, dbFile := range dbFiles {
		edb, err := gorm.Open(sqlite.Open(filepath.Join(config.BaseDir, dbFile)), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		var userAPIKeyED []UserAPIKeyED
		edb.Find(&userAPIKeyED)
		fmt.Printf("Found %d users", len(userAPIKeyED))

		var fmded []FileMetadataED
		edb.Find(&fmded)
		fmt.Printf("Found %d files", len(fmded))
	}
	var userAPiKeys []storagemodels.UserAPIKey
	var fmds []storagemodels.FileMetadata

	db.StorageDB.Find(&userAPiKeys)
	fmt.Printf("found %d users in new db", len(userAPiKeys))
	db.StorageDB.Find(&fmds)
	fmt.Printf("found %d files in new db", len(fmds))

}
func Migrate() {
	// check existing db, see if local db file name and exsting db same
	files, err := os.ReadDir(config.BaseDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		extension := filepath.Ext(file.Name())
		if (extension == ".db" || extension == ".sqlite") && file.Name() != config.StorageDbFile {
			edb, err := gorm.Open(sqlite.Open(filepath.Join(config.BaseDir, file.Name())), &gorm.Config{})
			if err != nil {
				panic(err)
			}
			var userAPIKeyED []UserAPIKeyED
			edb.Find(&userAPIKeyED)
			for _, userAPIKey := range userAPIKeyED {
				uakL := storagemodels.UserAPIKey{
					ID:       uuid.New(),
					Username: userAPIKey.Username,
					APIKey:   userAPIKey.APIKey,
					PIN:      userAPIKey.PIN,
				}
				db.StorageDB.Create(&uakL)
				logging.AuditLogger.Println("Created User ", userAPIKey.Username)
			}
			var fileMetadataED []FileMetadataED
			edb.Find(&fileMetadataED)
			for _, fmded := range fileMetadataED {
				fpSplit := strings.Split(fmded.FilePath, "/")
				username := fpSplit[0]
				user, err := storagemodels.GetUserByUsername(username)
				if err != nil {
					logging.Errorlogger.Println(username, fmded.FilePath, err)
				}
				fmdL := storagemodels.FileMetadata{
					Name:                       fmded.Name,
					FilePath:                   fmded.FilePath,
					SizeInMb:                   fmded.SizeInMb,
					IsPublic:                   fmded.IsPublic,
					IsImage:                    fmded.IsImage,
					CompressedVersionAvailable: fmded.CompressedVersionAvailable,
					CreatedAt:                  fmded.CreatedAt,
					UpdatedAt:                  fmded.UpdatedAt,
					UserID:                     user.ID,
				}
				logging.AuditLogger.Println("Created File ", fmdL.Name)
				db.StorageDB.Create(&fmdL)
			}

		}
	}

	// db := db.GetDB()
	// db.AutoMigrate(&storagemodels.File{}, &storagemodels.Folder{})
}
