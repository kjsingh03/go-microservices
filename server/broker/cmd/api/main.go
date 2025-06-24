package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct{}

func main() {

	_ = godotenv.Load("../.env") 										// ignore error - while running locally

	BROKER_PORT := os.Getenv("BROKER_PORT")								// While running via Docker Compose

	if BROKER_PORT == "" {												// Fallback
		BROKER_PORT = "8084" 
	}

	app := Config{}

	log.Printf("Broker Service starting at: http://localhost:%s", BROKER_PORT)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", BROKER_PORT),
		Handler: app.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
