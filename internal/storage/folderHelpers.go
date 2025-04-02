package storage

import (
	"fileupbackendv2/config"
	"os"
	"path/filepath"
)

func createFolder(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	} else if info.IsDir() {
		return nil
	} else {
		return err
	}
	return nil
}

func CreateFolder(username, subFolder string) error {
	var dirPath string
	if subFolder != "" {
		dirPath = filepath.Join(config.BaseDir, config.UploadsDir, username, subFolder)
	} else {
		dirPath = filepath.Join(config.BaseDir, config.UploadsDir, username)
	}
	return createFolder(dirPath)
}
