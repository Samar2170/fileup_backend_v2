package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var MODE string
var StorageBaseUrl string
var StorageBucketName string
var StorageAccessKey string
var BaseDir string

var ConnStr string
var LocalDBFile string
var SecretKey string

var LogsDir string

var APIPort string
var FrontEndPort string
var APIBaseUrl string

var StorageDbFile string

const (
	HeaderName         = "AccessKey"
	StorageSubDir      = "ecomm/products/"
	UploadsDir         = "uploads/"
	CompressionQuality = 85
)

func init() {
	currentFile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.Dir(currentFile)

	BaseDir = filepath.Dir(currentFile)
	// BaseDir = "/Users/samararora/Desktop/fileup-backend/"
	godotenv.Load(BaseDir + "/.env")

	MODE = os.Getenv("MODE")
	StorageAccessKey = os.Getenv("STORAGE_ACCESS_KEY")

	SecretKey = os.Getenv("SECRET_KEY")

	LogsDir = os.Getenv("LOGS_DIR")

	APIPort = os.Getenv("API_PORT")
	FrontEndPort = os.Getenv("FRONT_END_PORT")
	APIBaseUrl = os.Getenv("API_BASE_URL")

	StorageDbFile = os.Getenv("STORAGE_DB_FILE")
}
