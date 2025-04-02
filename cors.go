package main

import (
	"fileupbackendv2/config"
	"fileupbackendv2/pkg/logging"

	"github.com/rs/cors"
)

var CorsConfig *cors.Cors

func init() {
	allowedOrigins := []string{"https://backend.barebasics.shop/", "https://theclothingcompany.co", "https://storage.barebasics.shop/"}
	if config.MODE == "dev" {
		allowedOrigins = append(allowedOrigins, "http://localhost:9090")
	}
	logger := cors.Logger(&logging.AuditLogger)
	CorsConfig = cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Logger:           logger,
	})

}
