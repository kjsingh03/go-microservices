package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Manager struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewManager(uri, dbName string) (*Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Manager{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (m *Manager) GetDatabase() *mongo.Database {
	return m.db
}

func (m *Manager) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}