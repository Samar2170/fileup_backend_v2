package db

import (
	"fileupbackendv2/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var StorageDB *gorm.DB

func connectToStorageDB() {
	var err error
	StorageDB, err = gorm.Open(sqlite.Open(config.StorageDbFile), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func init() {
	connectToStorageDB()
}
