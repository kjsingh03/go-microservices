package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("../.env")

	if err!=nil{
		log.Panic("Failed to load env file")
	}
}

func GetEnvVar(key string,fallback string) string{
	key = os.Getenv(key)

	if key==""{
		key = fallback
	}

	return key
}