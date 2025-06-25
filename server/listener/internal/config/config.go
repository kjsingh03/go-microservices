package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQUser string
	RabbitMQPass string
	RabbitMQHost string
	RabbitMQPort string
}

func Load() (*Config) {
	_ = godotenv.Load("../.env")

	cfg := &Config{
		RabbitMQUser: getEnvVar("RABBITMQ_USER", "guest"),
		RabbitMQPass: getEnvVar("RABBITMQ_PASS", "guest"),
		RabbitMQHost: getEnvVar("RABBITMQ_HOST", "rabbitmq"),
		RabbitMQPort: getEnvVar("RABBITMQ_PORT", "5672"),
	}

	return cfg
}

func getEnvVar(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
