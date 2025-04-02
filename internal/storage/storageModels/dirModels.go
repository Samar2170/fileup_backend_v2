package storagemodels

import (
	"fileupbackendv2/internal/db"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Directory struct {
	*gorm.Model
	Name      string
	Path      string
	UserID    uuid.UUID
	User      UserAPIKey
	SizeInMb  float64
	CreatedAt time.Time
	UpdatedAt time.Time

	LastError   string
	HasError    bool `gorm:"default:false"`
	IsMasterDir bool `gorm:"default:false"`
}

func GetOrCreateDir(userId uuid.UUID, name string, isMasterDir bool) Directory {
	var dir Directory
	db.StorageDB.FirstOrCreate(&dir, Directory{UserID: userId, Name: name, IsMasterDir: isMasterDir})
	return dir
}

func GetDirByPathorName(path, name, username string) (Directory, error) {
	var dir Directory
	err := db.StorageDB.Where("name = ? OR path = ?", name, path).Where("username = ? ", username).Find(&dir).Error
	return dir, err

}
