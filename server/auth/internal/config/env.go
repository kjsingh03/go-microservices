package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}
}

func GetPort(key, fallback string) string {
	port := os.Getenv(key)
	if port == "" {
		return fallback
	}
	return port
}
