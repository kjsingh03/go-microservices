// internal/models/log.go
package models

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// LogEntry represents a log entry in the database
type LogEntry struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Data      string        `bson:"data" json:"data"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

// LogFilter can be used for filtering logs (future enhancement)
type LogFilter struct {
	Name      string    `json:"name,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Limit     int64     `json:"limit,omitempty"`
	Offset    int64     `json:"offset,omitempty"`
}

// LogStats represents statistics about logs
type LogStats struct {
	TotalLogs int64     `json:"total_logs"`
	OldestLog time.Time `json:"oldest_log,omitempty"`
	NewestLog time.Time `json:"newest_log,omitempty"`
}