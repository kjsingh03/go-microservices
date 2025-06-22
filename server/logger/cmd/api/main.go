package main

import (
	"fmt"
	"log"
	"logger/internal/config"
	"logger/internal/db"
	"logger/internal/models"
	"logger/internal/router"
	"net/http"
)

func main() {
	config.LoadEnv()

	mongoClient, err := db.ConnectToDB()

	if err != nil {
		log.Print(err)
	}

	models.InitModels(mongoClient)

	LOGGER_PORT := config.GetEnvVar("LOGGER_PORT", "")

	log.Printf("Server is starting at http://localhost:%s", LOGGER_PORT)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", LOGGER_PORT),
		Handler: router.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
