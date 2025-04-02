package main

import (
	"bufio"
	"context"
	"fileupbackendv2/config"
	"fileupbackendv2/frontend"
	"fileupbackendv2/handlers"
	"fileupbackendv2/internal/middleware"
	"fileupbackendv2/internal/storage"
	"fileupbackendv2/internal/utils"
	"fileupbackendv2/pkg/logging"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func checkKey() {
	fmt.Println("test")
	ak := getUserInput("Enter key:")
	fmt.Println(ak)
	hashed := utils.HashKey(ak)
	fmt.Println(hashed)
}
func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "setup":
			checkSetup()
		case "server":
			RunServer()
		case "new-user":
			generateKeyForUser()
		case "new-key":
			regenerateKeyForUser()
		case "check-key":
			checkKey()
		case "migrate":
			migrate()
		case "check-migration":
			checkMigration()
		case "dir-utils":
			storage.UpdateDirsData()
			storage.UpdateUserDirsData()
		case "feature":
			log.Println("Feature under development")
		case "help", "usage", "-h":
			fmt.Println("Usage: fileupbackendv2 [setup|server|new-user|regenerate-key|check-key|feature]")
		default:
			log.Fatal("Invalid argument")
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			RunServer()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			StartCronServer()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			frontend.StartEchoServer()
			wg.Done()
		}()
		wg.Wait()

	}
}

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	value = strings.ReplaceAll(value, "\n", "")
	return value
}

func generateKeyForUser() {
	username := getUserInput("Enter username: ")
	pin := getUserInput("Enter pin: ")
	ak, err := storage.CreateUser(username, pin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API Key: %s\n", ak)
}

func regenerateKeyForUser() {
	username := getUserInput("Enter username: ")
	pin := getUserInput("Enter pin: ")
	ak, err := storage.ReGenerateAPIKey(username, pin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API Key: %s\n", ak)
}

func checkSetup() {
	// check if dir exists
	uploadsDir := config.UploadsDir
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	logsDir := config.LogsDir
	if logsDir == "" {
		logsDir = config.BaseDir + "/logs"
	}

	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err = os.MkdirAll(logsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetServer() *http.Server {
	logger := logging.AuditLogger
	mux := mux.NewRouter()
	healthCheckHandler := http.HandlerFunc(handlers.HealthCheck)
	getFilesHandler := http.HandlerFunc(handlers.GetFilesHandler)
	uploadFileHandler := http.HandlerFunc(handlers.UploadFileHandler)
	getSignedUrlHandler := http.HandlerFunc(handlers.GetSignedUrlHandler)
	downloadHandler := http.HandlerFunc(handlers.DownloadFileHandler)

	mux.HandleFunc("/health-check/", healthCheckHandler)
	mux.HandleFunc("/files/get/", getFilesHandler)
	mux.HandleFunc("/files/upload/", uploadFileHandler).Methods("POST")
	mux.HandleFunc("/files/get-signed-url/{filepath:.*}", getSignedUrlHandler)
	mux.HandleFunc("/files/download/{filepath:.*}", downloadHandler)
	mux.HandleFunc("/folder/add/", handlers.CreateFolderHandler).Methods("POST")

	logMiddleware := logging.NewLogMiddleware(&logger)
	mux.Use(logMiddleware.Func())

	wrappedMux := middleware.APIKeyMiddleware(mux)
	wrappedMux = CorsConfig.Handler(wrappedMux)
	server := http.Server{
		Handler: wrappedMux,
		Addr:    ":" + config.APIPort,
	}
	return &server
}

func RunServer() {
	godotenv.Load()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server := GetServer()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("Running on %s\n", server.Addr)

	<-stop
	log.Println("Shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}
