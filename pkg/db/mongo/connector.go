package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ____________
// MongoDB
// ____________

type MongoDBConfig struct {
	Host string
	Port int
}

func LoadMongoDBConfig() (*MongoDBConfig, error) {
	host := os.Getenv("MONGODB_NAME")
	port, err := strconv.Atoi(os.Getenv("MONGODB_PORT"))
	if err != nil {
		return nil, errors.New("error: DB_PORT is not int")
	}
	return &MongoDBConfig{
		Host: host,
		Port: port,
	}, nil
}

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
		return errors.New("MongoDB already exists!")
	}

	// Create connection
	ctx := context.TODO()
	connectionStr := fmt.Sprintf("mongodb://%v:%v", cfg.Host, cfg.Port)
	log.Printf("Run mongodb server at %v", connectionStr)
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
