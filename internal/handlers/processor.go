package handlers

import (
	"context"
	"fmt"
	"log"

	myFMTer "github.com/https-whoyan/MafiaBot/internal/fmt"
	"github.com/https-whoyan/MafiaBot/internal/handlers/names"
	myTime "github.com/https-whoyan/MafiaBot/internal/time"
	"github.com/https-whoyan/MafiaBot/internal/util"
	"github.com/https-whoyan/MafiaBot/pkg"
	"github.com/https-whoyan/MafiaCore/game"

	"github.com/bwmarrin/discordgo"
)

type gameProcessor struct {
	g          *game.Game
	db         *pkg.Database
	f          *myFMTer.DiscordFMTer
	s          *discordgo.Session
	errCh      <-chan game.ErrSignal
	infoCh     <-chan game.InfoSignal
	errLogger  *log.Logger
	infoLogger *log.Logger
}

func newGameProcessor(
	g *game.Game, db *pkg.Database, s *discordgo.Session,
	errCh <-chan game.ErrSignal, infoCh <-chan game.InfoSignal,
	errLogger *log.Logger, infoLogger *log.Logger,
) *gameProcessor {
	return &gameProcessor{
		g:      g,
		db:     db,
		f:      myFMTer.DiscordFMTInstance,
		s:      s,
		errCh:  errCh,
		infoCh: infoCh,
		errLogger: log.New(
			errLogger.Writer(),
			errLogger.Prefix()+"gameProcessor; GuildID "+g.GuildID(),
			errLogger.Flags(),
		),
		infoLogger: log.New(
			infoLogger.Writer(),
			infoLogger.Prefix()+"gameProcessor; GuildID "+g.GuildID(),
			infoLogger.Flags(),
		),
	}
}

func (p *gameProcessor) process(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				break
			case err := <-p.errCh:
				p.errLogger.Printf("time: %v, error: %v", err.InitialTime, err.Err)
			case info := <-p.infoCh:
				p.infoLogger.Printf("got signal. time: %v, info: %v", info.InitialTime, info.InfoSignalType)
				if info.InfoSignalType != game.FinishGameSignal {
					return
				}
				p.offerRenameGame(ctx)
				return
			}
		}
	}(ctx)
}

func (p *gameProcessor) formatMessage(temporaryId string) string {
	return fmt.Sprintf("If you want to name a previous game, please and use " +
		"the " + names.RenameGameCommandName + " command and specify the " +
		p.f.Bl(temporaryId) + " identifier",
	)
}

// offerRenameGame
// This feature is needed to prompt users to reinitialize the game after completion
// Temporary game index generated in this function.
func (p *gameProcessor) offerRenameGame(ctx context.Context) {
	nameOfGame := p.g.GetStartTime().Format(myTime.BotTimeFormat)
	simplifyGame := game.DeepCloneGame{
		GuildID:   p.g.GuildID(),
		TimeStart: p.g.GetStartTime(),
	}
	_ = p.db.Storage.NameAGame(ctx, simplifyGame, nameOfGame)

	temporaryId := util.ToStr(util.GetRandomNumber())
	// Save to redis
	err := p.db.Hasher.SaveGameIndicator(ctx, temporaryId, simplifyGame)
	if err != nil {
		p.errLogger.Printf("could not save game: %v", err)
		return
	}
	_, err = sendMessages(p.s, p.g.GetMainChannel().GetServerID())
	if err != nil {
		p.errLogger.Printf("sendMessages: %v", err)
	}
}
