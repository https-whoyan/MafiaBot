package config

import (
	"log"
	"sync"

	"github.com/https-whoyan/MafiaBot/internal"
	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Config struct {
	MongoConfig *mongo.StorageConfig
	RedisConfig *redis.HasherConfig
	BotConfig   *bot.Config
}

var (
	configOnce sync.Once
)

func LoadConfig() *Config {
	var ansConfig *Config
	configOnce.Do(func() {
		loadDotEnv()

		mongoDBConfig, err := mongo.LoadMongoDBConfig()
		if err != nil {
			log.Fatal(err)
		}
		redisDBConfig, err := redis.LoadHasherConfig()
		if err != nil {
			log.Fatal(err)
		}
		botConfig := bot.LoadBotConfig()
		ansConfig = &Config{
			MongoConfig: mongoDBConfig,
			RedisConfig: redisDBConfig,
			BotConfig:   botConfig,
		}
	})
	if ansConfig == nil {
		log.Fatal("Config was been loaded before!")
	}
	log.Println("Load config...")
	return ansConfig
}

func logAboutDiscordGo() {
	log.Println("Discord-go version:", discordgo.VERSION)
	log.Println("Discord-go API version:", discordgo.APIVersion)
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
