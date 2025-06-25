package service

import (
	"context"
	"encoding/json"
	"fmt"
	"service-broker/internal/event"
	"service-broker/types"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitService struct {
	conn     *amqp.Connection
	emitter  *event.EventEmitter
	mu       sync.RWMutex
	closed   bool
}

// NewRabbitService creates a new RabbitMQ service
func NewRabbitService(conn *amqp.Connection) (RabbitService, error) {
	emitter, err := event.NewEventEmitter(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create event emitter: %w", err)
	}

	return &rabbitService{
		conn:    conn,
		emitter: emitter,
	}, nil
}

func (s *rabbitService) PublishLog(ctx context.Context, payload types.LogPayload) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("rabbit service is closed")
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal log payload: %w", err)
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := s.emitter.Push(string(jsonData), "log.INFO"); err != nil {
		return fmt.Errorf("failed to publish to queue: %w", err)
	}

	return nil
}

func (s *rabbitService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}