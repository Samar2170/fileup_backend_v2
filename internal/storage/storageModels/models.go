package storagemodels

import (
	"fileupbackendv2/internal/db"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	db.StorageDB.AutoMigrate(&UserAPIKey{})
	db.StorageDB.AutoMigrate(&FileMetadata{})
	db.StorageDB.AutoMigrate(&Directory{})
}

type UserAPIKey struct {
	*gorm.Model
	ID       uuid.UUID
	Username string
	APIKey   string
	PIN      string `gorm:"column:pin;varchar(6)"`
}

func (UserAPIKey) TableName() string {
	return "user_api_keys"
}

func GetUserById(id uint) (UserAPIKey, error) {
	var user UserAPIKey
	err := db.StorageDB.Where("id = ?", id).First(&user).Error
	return user, err
}
func GetUserByUsername(username string) (UserAPIKey, error) {
	var user UserAPIKey
	err := db.StorageDB.Where("username = ?", username).First(&user).Error
	return user, err
}

type FileMetadata struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"index"`
	FilePath string    `gorm:"index"`
	UserID   uuid.UUID `gorm:"index"`
	User     UserAPIKey
	SizeInMb float64
	IsPublic bool

	IsImage                    bool
	CompressedVersionAvailable bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
