package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL  string
	Name string
}

func Load() (*Config, error) {
	_ = godotenv.Load("../.env")

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnvWithDefault("LOGGER_PORT", "8085"),
			ReadTimeout:  getDurationWithDefault("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationWithDefault("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationWithDefault("IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			URL:  getEnvWithDefault("MONGO_URL", "mongodb://localhost:27017"),
			Name: getEnvWithDefault("DB_NAME", "logs"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}