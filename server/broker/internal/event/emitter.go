package event

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventEmitter struct {
	connection *amqp.Connection
}

func NewEventEmitter(conn *amqp.Connection) (*EventEmitter, error) {
	emitter := &EventEmitter{
		connection: conn,
	}

	// Setup exchange and queues
	if err := emitter.setup(); err != nil {
		return nil, fmt.Errorf("failed to setup emitter: %w", err)
	}

	return emitter, nil
}

func (e *EventEmitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	return declareExchange(ch)
}

func (e *EventEmitter) Push(body, routingKey string) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}

	return ch.Publish(
		"logs_topic", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		msg,          // message
	)
}