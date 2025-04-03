package auth

import (
	"errors"
	"fileupbackendv2/internal/db"
	"fileupbackendv2/internal/dirManager"
	"fileupbackendv2/internal/models"
	"fileupbackendv2/internal/utils"

	"github.com/google/uuid"
)

func CreateUser(username, email, password string) error {
	var user models.User
	var err error
	if err != nil {
		return err
	}
	usernameExists := models.CheckUsernameExists(username)
	if usernameExists {
		return errors.New("username already exists")
	}
	emailExists := models.CheckEmailExists(email)
	if emailExists {
		return errors.New("email already exists")
	}
	tx := db.StorageDB.Begin()
	user.ID = uuid.New()
	user.Username = username
	user.Email = email
	user.Password = utils.HashKey(password)
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dirManager.CreateFolder(username, "")
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func LoginUser(username, email, password string) (string, error) {
	var user models.User
	var err error
	if email == "" {
		err = db.StorageDB.Where("username = ?", username).Where("password = ?", utils.HashKey(password)).First(&user).Error
	} else {
		err = db.StorageDB.Where("email = ?", email).Where("password = ?", utils.HashKey(password)).First(&user).Error
	}
	token, err := utils.CreateToken(user.Username, user.ID.String())
	if err != nil {
		return "", err
	}
	return token, nil
}
