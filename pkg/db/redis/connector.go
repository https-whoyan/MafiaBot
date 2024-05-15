package redis

import "log"

type RedisConfig struct {
	//TODO!
}

func LoadRedisConfig() (*RedisConfig, error) {
	return &RedisConfig{}, nil
}

func InitRedis(cfg *RedisConfig) error {
	//TODO
	return nil
}

func Disconnect() error {
	//TODO!
	log.Println("Disconnect Redis")
	return nil
}
