package main

import (
	"listener/internal/config"
	"listener/internal/event"
	"listener/internal/rabbitmq"
	"log"
)

func main() {
	cfg := config.Load()

	rabbitConn, err := rabbitmq.Connect(cfg)
	if err != nil {
		log.Panic("RabbitMQ connection error:", err)
	}
	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}
