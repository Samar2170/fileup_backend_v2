package models

import (
	"fileupbackendv2/internal/db"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	db.StorageDB.AutoMigrate(&User{})
	db.StorageDB.AutoMigrate(&FileMetadata{})
	db.StorageDB.AutoMigrate(&Directory{})
}

type User struct {
	*gorm.Model
	ID         uuid.UUID
	Username   string `gorm:"index;unique"`
	Password   string
	Email      string `gorm:"index;unique"`
	APIKey     string
	IsVerified bool
}

func (User) TableName() string {
	return "users"
}
func CheckUsernameExists(username string) bool {
	var count int64
	db.StorageDB.Where("username = ?", username).Count(&count)
	return count > 0
}
func CheckEmailExists(email string) bool {
	var count int64
	db.StorageDB.Where("email = ?", email).Count(&count)
	return count > 0
}

func GetUserById(id string) (User, error) {
	var user User
	err := db.StorageDB.Where("id = ?", id).First(&user).Error
	return user, err
}
func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.StorageDB.Where("username = ?", username).First(&user).Error
	return user, err
}

type FileMetadata struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"index"`
	FilePath string    `gorm:"index"`
	UserID   uuid.UUID `gorm:"index"`
	User     User
	SizeInMb float64
	IsPublic bool

	IsImage                    bool
	CompressedVersionAvailable bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
