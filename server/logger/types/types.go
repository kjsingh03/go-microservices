// types/log.go
package types

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JsonResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Log struct {
	ID        string      `json:"id" bson:"_id,omitempty"`
	Name      string      `json:"name" bson:"name"`
	Data      interface{} `json:"data" bson:"data"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

type LogStats struct {
	TotalLogs    int64     `json:"total_logs"`
	LastLogTime  time.Time `json:"last_log_time"`
	FirstLogTime time.Time `json:"first_log_time"`
}

type CreateLogRequest struct {
	Name string      `json:"name" validate:"required"`
	Data interface{} `json:"data" validate:"required"`
}

type UpdateLogRequest struct {
	Name string      `json:"name,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// Custom BSON marshaling to handle ObjectID conversion
func (l *Log) MarshalBSON() ([]byte, error) {
	type Alias Log
	temp := &struct {
		ID interface{} `bson:"_id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	
	if l.ID != "" {
		if oid, err := bson.ObjectIDFromHex(l.ID); err == nil {
			temp.ID = oid
		} else {
			temp.ID = l.ID
		}
	}
	
	return bson.Marshal(temp)
}

// Custom BSON unmarshaling to handle ObjectID conversion
func (l *Log) UnmarshalBSON(data []byte) error {
	type Alias Log
	temp := &struct {
		ID interface{} `bson:"_id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	
	if err := bson.Unmarshal(data, temp); err != nil {
		return err
	}
	
	if temp.ID != nil {
		if oid, ok := temp.ID.(bson.ObjectID); ok {
			l.ID = oid.Hex()
		} else if str, ok := temp.ID.(string); ok {
			l.ID = str
		}
	}
	
	return nil
}