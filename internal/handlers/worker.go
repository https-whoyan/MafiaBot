package handlers

import (
	"context"
	"fmt"
	"github.com/https-whoyan/MafiaBot/internal/handlers/names"
	"github.com/https-whoyan/MafiaBot/internal/util"
	"time"

	bot "github.com/https-whoyan/MafiaBot/internal"
	myTime "github.com/https-whoyan/MafiaBot/internal/time"
	"github.com/https-whoyan/MafiaCore/game"
)

type gameCleaner struct {
	bot *bot.Bot
	g   *game.Game
}

const clearTiming = 10 * time.Second

func NewGameCleaner(bot *bot.Bot, g *game.Game) *gameCleaner {
	return &gameCleaner{
		bot: bot,
		g:   g,
	}
}

func (c *gameCleaner) Work(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(clearTiming):
				c.check(ctx)
			}
		}
	}(ctx)
}

func (c *gameCleaner) formatMessage(identificator string) string {
	return fmt.Sprintf("If you want to name a previous game, please and use " +
		"the " + names.RenameGameCommandName + " command and specify the " +
		c.bot.FMTer.Bl(identificator) + " identifier")
}

func (c *gameCleaner) check(ctx context.Context) {
	if c.g.GetState() != game.FinishState {
		return
	}
	// name of game
	nameOfGame := c.g.GetStartTime().Format(myTime.BotTimeFormat)
	simplyfyGame := game.DeepCloneGame{
		GuildID:   c.g.GuildID(),
		TimeStart: c.g.GetStartTime(),
	}
	_ = c.g.GetLogger().NameAGame(ctx, simplyfyGame, nameOfGame)

	identificator := util.ToStr(util.GetRandomNumber())
	// Save to redis
	err := c.bot.Databases.Hasher.SaveGameIndicator(ctx, identificator, simplyfyGame)
	if err != nil {
		return
	}
	message := c.formatMessage(identificator)
	err = c.g.GameMessenger().Public.SendMessageToMainChat(message)
	if err != nil {
		return
	}
}
