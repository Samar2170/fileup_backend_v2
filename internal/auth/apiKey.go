package auth

import (
	"fileupbackendv2/internal/db"
	"fileupbackendv2/internal/models"
	"fileupbackendv2/internal/utils"
)

func GenerateAPIKey(username string, password string) (string, error) {
	var user models.User
	newKey, err := utils.GenerateKey(32)
	if err != nil {
		return "", err
	}
	hashedPassword := utils.HashKey(password)
	err = db.StorageDB.Where("username = ?", username).Where("password = ?", hashedPassword).First(&user).Error
	if err != nil {
		return "", err
	}
	user.APIKey = utils.HashKey(newKey)
	db.StorageDB.Save(&user)
	return newKey, nil
}

func IsKeyValid(key string) bool {
	keyHash := utils.HashKey(key)
	var apiKey models.User
	err := db.StorageDB.Where("api_key = ?", keyHash).First(&apiKey).Error
	return err == nil
}

func GetUserByKey(key string) (models.User, error) {
	keyHash := utils.HashKey(key)
	var apiKey models.User
	err := db.StorageDB.Where("api_key = ?", keyHash).First(&apiKey).Error
	if err != nil {
		return apiKey, err
	}
	return apiKey, nil
}
