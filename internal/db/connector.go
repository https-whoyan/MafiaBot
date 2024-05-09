package db

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

var currDB *MongoDB

func InitMongoDB(cfg *MongoDBConfig) error {
	// Check is containing
	if currDB != nil {
		return errors.New("MongoDB already exists!")
	}

	// Create connection
	ctx := context.TODO()
	connectionStr := fmt.Sprintf("mongodb://%v:%v", cfg.Port, strconv.Itoa(cfg.Port))
	log.Printf("Run mongodb server at %v", connectionStr)
	clientOptions := options.Client().ApplyURI(connectionStr)
	client, err := mongo.Connect(ctx, clientOptions)
	fmt.Println(client)
	if err != nil {
		return err
	}

	// Check is ok
	err = client.Ping(ctx, nil)
	fmt.Println("я тут 3")
	if err != nil {
		return err
	}

	// Initial currDB singleton var
	currDB = &MongoDB{
		db: client,
	}

	// disconnect function
	defer func(db *MongoDB) {
		err = currDB.db.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}(currDB)
	return err
}

// GetCurrDB get connection
func GetCurrDB() (*MongoDB, bool) {
	if currDB == nil {
		return nil, false
	}
	return currDB, true
}
