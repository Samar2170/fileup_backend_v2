package storage

import (
	"fileupbackendv2/internal/db"
	storagemodels "fileupbackendv2/internal/storage/storageModels"
	"fileupbackendv2/internal/utils"

	"github.com/google/uuid"
)

func CreateUser(username, pin string) (string, error) {
	var user storagemodels.UserAPIKey
	newKey, err := utils.GenerateKey(32)
	if err != nil {
		return "", err
	}
	tx := db.StorageDB.Begin()
	user.ID = uuid.New()
	user.Username = username
	user.PIN = pin
	user.APIKey = utils.HashKey(newKey)
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = CreateFolder(username, "")
	if err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return newKey, nil
}

func ReGenerateAPIKey(username string, pin string) (string, error) {
	var user storagemodels.UserAPIKey
	newKey, err := utils.GenerateKey(32)
	if err != nil {
		return "", err
	}
	err = db.StorageDB.Where("username = ?", username).Where("pin = ?", pin).First(&user).Error
	if err != nil {
		return "", err
	}
	user.APIKey = utils.HashKey(newKey)
	db.StorageDB.Save(&user)
	return newKey, nil
}

func IsKeyValid(key string) bool {
	keyHash := utils.HashKey(key)
	var apiKey storagemodels.UserAPIKey
	err := db.StorageDB.Where("api_key = ?", keyHash).First(&apiKey).Error
	return err == nil
}

func GetUserByKey(key string) (storagemodels.UserAPIKey, error) {
	keyHash := utils.HashKey(key)
	var apiKey storagemodels.UserAPIKey
	err := db.StorageDB.Where("api_key = ?", keyHash).First(&apiKey).Error
	if err != nil {
		return apiKey, err
	}
	return apiKey, nil
}
