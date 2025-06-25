// internal/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Environment string
	Server      ServerConfig
	RabbitMQ    RabbitMQConfig
	Services    ServicesConfig
	Rabbit      *amqp.Connection
}

type ServerConfig struct {
	Port string
	Host string
}

type RabbitMQConfig struct {
	URL             string
	Host            string
	Port            string
	Username        string
	Password        string
	VHost           string
	ConnectionName  string
	ConnectionRetry int
}

type ServicesConfig struct {
	AuthURL    string
	LogURL     string
	MailURL    string
	Timeout    time.Duration
	RetryCount int
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		Environment: GetEnvVar("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port: GetEnvVar("BROKER_PORT", "80"),
			Host: GetEnvVar("BROKER_HOST", "0.0.0.0"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:            GetEnvVar("RABBITMQ_HOST", "localhost"),
			Port:            GetEnvVar("RABBITMQ_PORT", "5672"),
			Username:        GetEnvVar("RABBITMQ_USER", "guest"),
			Password:        GetEnvVar("RABBITMQ_PASS", "guest"),
			VHost:           GetEnvVar("RABBITMQ_VHOST", "/"),
			ConnectionName:  GetEnvVar("RABBITMQ_CONNECTION_NAME", "broker-service"),
			ConnectionRetry: GetEnvInt("RABBITMQ_CONNECTION_RETRY", 5),
		},
		Services: ServicesConfig{
			AuthURL:    GetEnvVar("AUTH_SERVICE_URL", "http://authentication-service"),
			LogURL:     GetEnvVar("LOG_SERVICE_URL", "http://logger-service/api/v1"),
			MailURL:    GetEnvVar("MAIL_SERVICE_URL", "http://mailer-service/api/v1"),
			Timeout:    GetEnvDuration("SERVICE_TIMEOUT", 30*time.Second),
			RetryCount: GetEnvInt("SERVICE_RETRYCOUNT", 5),
		},
	}

	cfg.RabbitMQ.URL = fmt.Sprintf("amqp://%s:%s@%s:%s%s", cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port, cfg.RabbitMQ.VHost)

	conn, err := ConnectRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	cfg.Rabbit = conn

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func ConnectRabbitMQ(cfg RabbitMQConfig) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < cfg.ConnectionRetry; i++ {
		conn, err = amqp.Dial(cfg.URL)
		if err == nil {
			log.Printf("Successfully connected to RabbitMQ at %s:%s", cfg.Host, cfg.Port)
			return conn, nil
		}

		if i < cfg.ConnectionRetry-1 {
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", cfg.ConnectionRetry, err)
}
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Services.AuthURL == "" {
		return fmt.Errorf("auth service URL is required")
	}
	if c.Services.LogURL == "" {
		return fmt.Errorf("log service URL is required")
	}
	if c.Services.MailURL == "" {
		return fmt.Errorf("mail service URL is required")
	}
	return nil
}

func GetEnvVar(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
