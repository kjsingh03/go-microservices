package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load("../.env")
}

func GetPort(key, fallback string) string {
	port := os.Getenv(key)
	if port == "" {
		return fallback
	}
	return port
}
