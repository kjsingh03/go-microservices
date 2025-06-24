package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	// Database DatabaseConfig
	Mailer   MailerConfig
}

type ServerConfig struct {
	Port    string
	Handler *http.Handler
}

// type DatabaseConfig struct {
// 	Name string
// 	Url  string
// }

type MailerConfig struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

func Load() (*Config, error) {
	_ = godotenv.Load("../.env")

	mailPort, err := strconv.Atoi(getEnvVar("MAIL_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAIL_PORT: %w", err)
	}

	app := &Config{
		Server: ServerConfig{
			Port: getEnvVar("MAILER_PORT", "8082"),
		},
		// Database: DatabaseConfig{
		// 	Name: getEnvVar("MAILER_DB_NAME", "mailer"),
		// 	Url:  getEnvVar("MONGO_URL", ""),
		// },
		Mailer: MailerConfig{
			Domain:      getEnvVar("MAIL_DOMAIN", ""),
			Host:        getEnvVar("MAIL_HOST", ""),
			Port:        mailPort,
			Username:    getEnvVar("MAIL_USERNAME", ""),
			Password:    getEnvVar("MAIL_PASSWORD", ""),
			Encryption:  strings.ToLower(getEnvVar("MAIL_ENCRYPTION", "tls")),
			FromName:    getEnvVar("FROM_NAME", ""),
			FromAddress: getEnvVar("FROM_ADDRESS", ""),
		},
	}

	if err := app.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return app, nil
}

func getEnvVar(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (c *Config) ValidateConfig() error {
	var errors []string

	// if c.Database.Url == "" {
	// 	errors = append(errors, "MONGO_URL is required")
	// }

	if c.Mailer.Host == "" {
		errors = append(errors, "MAIL_HOST is required")
	}
	if c.Mailer.Username == "" {
		errors = append(errors, "MAIL_USERNAME is required")
	}
	if c.Mailer.Password == "" {
		errors = append(errors, "MAIL_PASSWORD is required")
	}
	if c.Mailer.FromAddress == "" {
		errors = append(errors, "FROM_ADDRESS is required")
	}
	if c.Mailer.FromName == "" {
		errors = append(errors, "FROM_NAME is required")
	}

	validEncryptions := []string{"tls", "ssl", "none", ""}
	isValidEncryption := false
	for _, enc := range validEncryptions {
		if c.Mailer.Encryption == enc {
			isValidEncryption = true
			break
		}
	}
	if !isValidEncryption {
		errors = append(errors, "MAIL_ENCRYPTION must be one of: tls, ssl, none")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration errors: %s", strings.Join(errors, "; "))
	}

	return nil
}