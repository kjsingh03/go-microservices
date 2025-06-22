// models/models.go
package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Global models instance
var AppModels Models

type Models struct {
	Log LogEntryModel
}

type Log struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Data      string        `bson:"data" json:"data"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type LogEntryModel struct {
	Collection *mongo.Collection
}

func InitModels(client *mongo.Client) {
	AppModels = Models{
		Log: LogEntryModel{
			Collection: client.Database("logs").Collection("logs"),
		},
	}
}

func (m *LogEntryModel) Insert(entry Log) error {
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	_, err := m.Collection.InsertOne(context.Background(), entry)
	if err != nil {
		log.Println("Insert error:", err)
		return err
	}
	return nil
}

func (m *LogEntryModel) All() ([]*Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := m.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Find error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*Log
	for cursor.Next(ctx) {
		var logEntry Log
		if err := cursor.Decode(&logEntry); err != nil {
			log.Println("Decode error:", err)
			return nil, err
		}
		logs = append(logs, &logEntry)
	}
	return logs, nil
}

func (m *LogEntryModel) GetOne(id string) (*Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	docID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry Log
	err = m.Collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (m *LogEntryModel) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return m.Collection.Drop(ctx)
}

func (m *LogEntryModel) Update(entry *Log) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	docID := entry.ID
	entry.UpdatedAt = time.Now()

	update := bson.D{
		{"$set", bson.D{
			{"name", entry.Name},
			{"data", entry.Data},
			{"updated_at", entry.UpdatedAt},
		}},
	}

	return m.Collection.UpdateOne(ctx, bson.M{"_id": docID}, update)
}

func (m *LogEntryModel) Delete(id bson.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return m.Collection.DeleteOne(ctx, bson.M{"_id": id})
}