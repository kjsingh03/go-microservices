// internal/repositories/log_repository.go
package repositories

import (
	"context"
	"logger/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// logRepository implements LogRepository interface
type logRepository struct {
	collection *mongo.Collection
}

// NewLogRepository creates a new log repository instance
func NewLogRepository(db *mongo.Database) types.LogRepositoryInterface {
	return &logRepository{
		collection: db.Collection("logs"),
	}
}

func (r *logRepository) Create(ctx context.Context, log *types.Log) error {
	// Set timestamps
	now := time.Now().UTC()
	log.CreatedAt = now
	log.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	
	// Set the ID from MongoDB's generated ObjectID
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		log.ID = oid.Hex()
	}
	
	return nil
}

func (r *logRepository) FindAll(ctx context.Context) ([]types.Log, error) {
	// Sort by created_at in descending order (newest first)
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	
	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []types.Log
	for cursor.Next(ctx) {
		var log types.Log
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, cursor.Err()
}

func (r *logRepository) FindByID(ctx context.Context, id string) (*types.Log, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var log types.Log
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil when not found
		}
		return nil, err
	}

	// Set the ID field from the ObjectID
	log.ID = objectID.Hex()
	return &log, nil
}

func (r *logRepository) Update(ctx context.Context, id string, log *types.Log) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Update timestamp
	log.UpdatedAt = time.Now().UTC()

	filter := bson.M{"_id": objectID}
	update := bson.D{
		{"$set", bson.D{
			{"name", log.Name},
			{"data", log.Data},
			{"updated_at", log.UpdatedAt},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	
	return nil
}

func (r *logRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	
	return nil
}

func (r *logRepository) DropCollection(ctx context.Context) error {
	return r.collection.Drop(ctx)
}

func (r *logRepository) GetStats(ctx context.Context) (*types.LogStats, error) {
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return &types.LogStats{
		TotalLogs: count,
	}, nil
}