package db

import (
	"context"
	"fmt"
	"logger/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
) 

func ConnectToDB() (*mongo.Client, error) {
	
	MONGO_URL := config.GetEnvVar("MONGO_URL", "mongodb://localhost:27017")

	if MONGO_URL == "" {
		return nil, fmt.Errorf("failed to load Mongo URL")
	}

	clientOptions := options.Client().ApplyURI(MONGO_URL)

	// creds := options.Credential{
	// 	Username: config.GetEnvVar("MONGO_USER", ""),
	// 	Password: config.GetEnvVar("MONGO_PASS", ""),
	// }

	// clientOptions.SetAuth(creds)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	return client, nil
}
