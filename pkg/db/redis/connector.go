package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadRedisConfig() (*RedisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	return &RedisConfig{
		Host: host,
		Port: port,
		DB:   db,
	}, nil
}

type RedisDB struct {
	sync.Mutex
	db *redis.Client
}

var (
	redisOnce   sync.Once
	currRedisDB *RedisDB
)

func InitRedis(cfg *RedisConfig) error {
	// Create connection string
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

	// Background context
	ctx := context.Background()
	// Create client
	client := redis.NewClient(connectionOptions)

	// Check is ok
	status := client.Ping(ctx)
	val, err := status.Result()
	if err != nil {
		return err
	}

	if val != "PONG" {
		return errors.New(
			fmt.Sprintf(
				"excepted PONG, get %v", val,
			),
		)
	}

	redisOnce.Do(func() {
		currRedisDB = &RedisDB{
			db: client,
		}
	})
	log.Printf("Run redis server at %v, db: %v", connectionStr, cfg.DB)

	return nil
}

func GetCurrRedisDB() (*RedisDB, bool) {
	if currRedisDB == nil {
		return nil, false
	}
	return currRedisDB, true
}

func Disconnect() error {
	if currRedisDB == nil {
		return errors.New("redis DB isn't initialized")
	}
	currRedisDB.Lock()
	err := currRedisDB.db.Close()
	currRedisDB.Unlock()
	log.Println("Disconnect Redis")
	return err
}
