package rabbitmq

import (
	"fmt"
	"listener/internal/config"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(cfg *config.Config) (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQUser, cfg.RabbitMQPass, cfg.RabbitMQHost, cfg.RabbitMQPort)

	for {
		c, err := amqp.Dial(url)
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			log.Println("Too many retries, exiting.")
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off for", backOff)
		time.Sleep(backOff)
	}

	return connection, nil
}