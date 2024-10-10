package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

type HasherConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Logger   *log.Logger
}

func LoadHasherConfig() (*HasherConfig, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	return &HasherConfig{
		Host: host,
		Port: port,
		DB:   db,
	}, nil
}

func (c *HasherConfig) SetLogger(logger *log.Logger) {
	c.Logger = logger
}

func InitHasher(ctx context.Context, cfg *HasherConfig) (Hasher, error) {
	connectionStr := fmt.Sprintf(
		"%v:%v",
		cfg.Host,
		cfg.Port)
	// Create connection options
	connectionOptions := &redis.Options{
		Addr:     connectionStr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	// Create client
	client := redis.NewClient(connectionOptions)

	// Check is ok
	status := client.Ping(ctx)
	val, err := status.Result()
	if err != nil {
		return nil, err
	}

	if val != "PONG" {
		return nil, errors.New(
			fmt.Sprintf(
				"excepted PONG, get %v", val,
			),
		)
	}

	cfg.Logger.Printf("Run redis server at %v, db: %v", connectionStr, cfg.DB)
	return &redisDB{
		db: client,
		lg: cfg.Logger,
	}, nil
}
