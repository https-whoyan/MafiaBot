package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/https-whoyan/MafiaBot/internal"
	"github.com/https-whoyan/MafiaBot/internal/app/config"
	"github.com/https-whoyan/MafiaBot/pkg"
	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"
)

type gracefulShutdown[T any] struct {
	operand T
	close   func(ctx context.Context, operator T) error
}

func newGracefulShutdown[T any](operand T, closerFunc func(ctx context.Context, operand T) error) gracefulShutdown[T] {
	return gracefulShutdown[T]{
		operand: operand,
		close:   closerFunc,
	}
}

func (g gracefulShutdown[T]) listen(ctx context.Context) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	cancelSignal := <-ch
	err := g.close(ctx, g.operand)
	if err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Printf("graceful finding process error: %v", err)
		return
	}
	_ = p.Signal(cancelSignal)
}

type App struct {
	Bot       *bot.Bot
	Databases *pkg.Database
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
	discordBot, err := bot.InitBot(ctx, cfg.BotConfig, databases)
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

	hasherGracefulShutdown := newGracefulShutdown[redis.Hasher](a.Databases.Hasher, func(ctx context.Context, hasher redis.Hasher) error {
		return hasher.Close(ctx)
	})
	storageGracefulShutdown := newGracefulShutdown[mongo.Storage](a.Databases.Storage, func(ctx context.Context, storage mongo.Storage) error {
		return storage.Close(ctx)
	})
	botGracefulShutdown := newGracefulShutdown[*bot.Bot](a.Bot, func(ctx context.Context, bot *bot.Bot) error {
		return bot.Close()
	})
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go hasherGracefulShutdown.listen(ctx)
	go storageGracefulShutdown.listen(ctx)
	botGracefulShutdown.listen(ctx)
	return nil
}
