package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// ____________
// Config
// ____________

type StorageConfig struct {
	Host string
	Port string
}

func LoadMongoDBConfig() (*StorageConfig, error) {
	host := os.Getenv("MONGODB_NAME")
	port := os.Getenv("MONGODB_PORT")
	return &StorageConfig{
		Host: host,
		Port: port,
	}, nil
}

// ____________
// mongoDB
// ____________

type mongoDB struct {
	db *mongo.Client
}

func InitStorage(ctx context.Context, cfg *StorageConfig) (Storage, error) {
	connectionStr := fmt.Sprintf("mongodb://%v:%v", cfg.Host, cfg.Port)
	clientOptions := options.Client().ApplyURI(connectionStr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check is ok
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Run mongodb server at %v", connectionStr)
	return &mongoDB{client}, nil
}

func (s *mongoDB) Close(ctx context.Context) error {
	return s.db.Disconnect(ctx)
}
