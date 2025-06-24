package models

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type LogEntry struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Data      string        `bson:"data" json:"data"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type LogStats struct {
	TotalLogs int64     `json:"total_logs"`
	OldestLog time.Time `json:"oldest_log,omitempty"`
	NewestLog time.Time `json:"newest_log,omitempty"`
}