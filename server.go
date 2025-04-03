package main

import (
	"context"
	"fileupbackendv2/config"
	"fileupbackendv2/handlers"
	"fileupbackendv2/internal/middleware"
	"fileupbackendv2/pkg/logging"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func GetServer() *http.Server {
	logger := logging.AuditLogger
	mux := mux.NewRouter()

	mux.HandleFunc("/health-check/", handlers.HealthCheck)

	mux.HandleFunc("/auth/signup/", handlers.SignupHandler).Methods("POST")
	mux.HandleFunc("/auth/login/", handlers.LoginHandler).Methods("POST")
	mux.HandleFunc("/auth/generate-api-key/", handlers.GenerateAPIKeyHandler).Methods("POST")

	mux.HandleFunc("/files/get/", handlers.GetFilesHandler)
	mux.HandleFunc("/files/upload/", handlers.UploadFileHandler).Methods("POST")
	mux.HandleFunc("/files/get-signed-url/{filepath:.*}", handlers.GetSignedUrlHandler)
	mux.HandleFunc("/files/download/{filepath:.*}", handlers.DownloadFileHandler)
	mux.HandleFunc("/folder/add/", handlers.CreateFolderHandler).Methods("POST")

	logMiddleware := logging.NewLogMiddleware(&logger)
	mux.Use(logMiddleware.Func())

	wrappedMux := middleware.AuthMiddleware(mux)
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
	fmt.Printf("Running in %s mode\n", config.MODE)

	<-stop
	log.Println("Shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped")
}
