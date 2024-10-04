package workers

import (
	"context"
	bot "github.com/https-whoyan/MafiaBot/internal"
	"github.com/https-whoyan/MafiaCore/game"
	"time"
)

type GameCleaner struct {
	bot *bot.Bot
}

const cleanerDuration = 10 * time.Second

func NewGameCleaner(bot *bot.Bot) *GameCleaner {
	return &GameCleaner{
		bot: bot,
	}
}

func (c *GameCleaner) Work(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(cleanerDuration)
				c.ckeck()
			}
		}
	}(ctx)
}

func (c *GameCleaner) ckeck() {
	for guildID, guildGame := range c.bot.Games {
		if guildGame.GetState() == game.FinishState {
			_ = guildGame.GetLogger().NameAGame(ctx)
		}
	}
}
