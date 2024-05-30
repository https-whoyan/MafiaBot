package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ____________
// Config
// ____________

type MongoDBConfig struct {
	Host string
	Port string
}

func LoadMongoDBConfig() (*MongoDBConfig, error) {
	host := os.Getenv("MONGODB_NAME")
	port := os.Getenv("MONGODB_PORT")
	return &MongoDBConfig{
		Host: host,
		Port: port,
	}, nil
}

// ____________
// MongoDB
// ____________

type MongoDB struct {
	sync.Mutex
	db *mongo.Client
}

var (
	mongoOnce   sync.Once
	currMongoDB *MongoDB
)

func InitMongoDB(cfg *MongoDBConfig) error {
	// Check is containing
	if currMongoDB != nil {
		return errors.New("mongoDB already exists")
	}

	// Create connection
	ctx := context.TODO()
	connectionStr := fmt.Sprintf("mongodb://%v:%v", cfg.Host, cfg.Port)
	clientOptions := options.Client().ApplyURI(connectionStr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Check is ok
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Initial currDB singleton var
	mongoOnce.Do(func() {
		currMongoDB = &MongoDB{
			db: client,
		}
	})
	log.Printf("Run mongodb server at %v", connectionStr)

	return err
}

// GetCurrMongoDB get connection
func GetCurrMongoDB() (*MongoDB, bool) {
	if currMongoDB == nil {
		return nil, false
	}
	return currMongoDB, true
}

func DisconnectMongoDB() error {
	if currMongoDB == nil {
		return errors.New("mongoDB is not initialzed")
	}
	currMongoDB.Lock()
	err := currMongoDB.db.Disconnect(context.TODO())
	currMongoDB.Unlock()
	log.Println("Disconnect MongoDB")
	return err
}
