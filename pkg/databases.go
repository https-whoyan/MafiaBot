package pkg

import (
	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"
)

type Database struct {
	Storage mongo.Storage
	Hasher  redis.Hasher
}

func NewDatabase(storage mongo.Storage, hasher redis.Hasher) *Database {
	return &Database{
		Storage: storage,
		Hasher:  hasher,
	}
}
