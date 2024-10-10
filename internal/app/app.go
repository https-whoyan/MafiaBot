package app

import (
	"context"
	"log"

	bot "github.com/https-whoyan/MafiaBot/internal"
	"github.com/https-whoyan/MafiaBot/internal/app/config"
	"github.com/https-whoyan/MafiaBot/pkg"
	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"
)

type App struct {
	Bot         *bot.Bot
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Databases   *pkg.Database
}

func InitApp(ctx context.Context, cfg *config.Config) (*App, error) {
	mongoStorage, err := mongo.InitStorage(ctx, cfg.MongoConfig)
	if err != nil {
		return nil, err
	}
	redisStorage, err := redis.InitHasher(ctx, cfg.RedisConfig)
	if err != nil {
		return nil, err
	}
	databases := pkg.NewDatabase(mongoStorage, redisStorage)
	discordBot, err := bot.InitBot(
		ctx,
		cfg.BotConfig, databases, cfg.ErrorLogger, cfg.InfoLogger,
	)
	if err != nil {
		return nil, err
	}
	return &App{
		Bot:       discordBot,
		Databases: databases,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	// bot
	err := a.Bot.Open()
	if err != nil {
		return err
	}

	hasherGracefulShutdown := newGracefulShutdown[redis.Hasher](
		a.Databases.Hasher,
		func(ctx context.Context, hasher redis.Hasher) error {
			return hasher.Close(ctx)
		},
		a.ErrorLogger,
	)
	storageGracefulShutdown := newGracefulShutdown[mongo.Storage](
		a.Databases.Storage,
		func(ctx context.Context, storage mongo.Storage) error {
			return storage.Close(ctx)
		},
		a.ErrorLogger,
	)
	botGracefulShutdown := newGracefulShutdown[*bot.Bot](
		a.Bot,
		func(ctx context.Context, bot *bot.Bot) error {
			return bot.Close()
		},
		a.ErrorLogger,
	)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go hasherGracefulShutdown.listen(ctx)
	go storageGracefulShutdown.listen(ctx)
	botGracefulShutdown.listen(ctx)
	return nil
}
