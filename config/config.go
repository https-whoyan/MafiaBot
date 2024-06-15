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
	mongoConfig *mongo.MongoDBConfig
	redisConfig *redis.RedisConfig
	botConfig   *bot.BotConfig
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
		redisDBConfig, err := redis.LoadRedisConfig()
		if err != nil {
			log.Fatal(err)
		}
		botConfig := bot.LoadBotConfig()
		ansConfig = &Config{
			mongoConfig: mongoDBConfig,
			redisConfig: redisDBConfig,
			botConfig:   botConfig,
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

func (c *Config) Run() {
	err := mongo.InitMongoDB(c.mongoConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = redis.InitRedis(c.redisConfig)
	if err != nil {
		log.Fatal(err)
	}

	logAboutDiscordGo()
	bot.InitBot(c.botConfig)

	bot.Run()

	defer func() {
		err = redis.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
		err = mongo.DisconnectMongoDB()
		if err != nil {
			log.Fatal(err)
		}
		err = bot.DisconnectBot()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
