package redis

import (
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

type RedisConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       int
}

func LoadRedisConfig() (*RedisConfig, error) {
	return &RedisConfig{}, nil
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
	/*
		redisOnce.Do(func() {
			connectionStr := fmt.Sprintf(
				"%v:%v",
				cfg.Host,
				cfg.Port)
			//.Println(connectionStr)
		})
		//TODO!
	*/
	return nil
}

func Disconnect() error {
	//TODO!
	log.Println("Disconnect Redis")
	return nil
}
