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
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

var (
	configOnce sync.Once
)

func LoadConfig() *Config {
	log.Println("Load config...")
	var ansConfig *Config
	configOnce.Do(func() {
		loadDotEnv()
		infoLogger, errLogger := initLoggers()
		mongoDBConfig, err := mongo.LoadMongoDBConfig()
		if err != nil {
			log.Fatal(err)
		}
		redisDBConfig, err := redis.LoadHasherConfig()
		if err != nil {
			log.Fatal(err)
		}
		botConfig := bot.LoadBotConfig()

		mongoDBConfig.SetLogger(infoLogger)
		redisDBConfig.SetLogger(infoLogger)
		ansConfig = &Config{
			MongoConfig: mongoDBConfig,
			RedisConfig: redisDBConfig,
			BotConfig:   botConfig,
			InfoLogger:  infoLogger,
			ErrorLogger: errLogger,
		}
		ansConfig.infoAboutDiscordGo()
	})
	if ansConfig == nil {
		log.Fatal("Config was been loaded before!")
	}
	return ansConfig
}

func (c *Config) infoAboutDiscordGo() {
	c.InfoLogger.Println("Discord-go version:", discordgo.VERSION)
	c.InfoLogger.Println("Discord-go API version:", discordgo.APIVersion)
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
