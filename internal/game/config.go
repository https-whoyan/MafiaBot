package game

import (
	"log"
	"math/rand"
	"sort"

	botFMT "github.com/https-whoyan/MafiaBot/internal/fmt"
	coreUserPack "github.com/https-whoyan/MafiaBot/internal/user"

	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"

	coreGamePack "github.com/https-whoyan/MafiaCore/game"
)

func GetNewGameConfig(
	renameProvider *coreUserPack.BotUserRenameProvider, storage mongo.Storage,
	errLogger *log.Logger, infoLogger *log.Logger,
) []coreGamePack.Option {
	options := []coreGamePack.Option{
		coreGamePack.FMTerOpt(botFMT.DiscordFMTInstance),
		coreGamePack.RenameModeOpt(coreGamePack.RenameInGuildMode),
		coreGamePack.RenamePrOpt(renameProvider),
		coreGamePack.VotePingOpt(1),
		coreGamePack.StorageOpt(storage),
		coreGamePack.ErrLoggerOpt(errLogger),
		coreGamePack.InfoLoggerOpt(infoLogger),
	}
	return options
}

// ConfigMessages A structure that stores information, which Redis will use to store information about all
// messages where people should put reactions to the configuration they like.
//
// This type declarative here to avoid recursive import.
type ConfigMessages struct {
	// So that after voting, the bot can quickly index which configurations were sent.
	PlayersCount int `redis:"playersCount"`
	GuildID      string
	ConfigsCount int              `redis:"configsCount"`
	Messages     []*configMessage `redis:"messages"`
}

type configMessage struct {
	ConfigIndex    int    `redis:"configIndex"`
	MessageID      string `redis:"messageID"`
	ReactionsCount int
}

func NewConfigMessages(guildID string, peoplesCount int, configsCount int) *ConfigMessages {
	return &ConfigMessages{
		PlayersCount: peoplesCount,
		GuildID:      guildID,
		ConfigsCount: configsCount,
		Messages:     make([]*configMessage, 0),
	}
}

// Setters

func (c *configMessage) SetReactionCount(reactionsCount int) { c.ReactionsCount = reactionsCount }
func (c *ConfigMessages) AddNewMessage(configIndex int, messageID string) {
	c.Messages = append(c.Messages, &configMessage{
		ConfigIndex: configIndex,
		MessageID:   messageID,
	})
}

// Utils

func (c *ConfigMessages) GetWinner() (configIndex int, playersCount int, isRandom bool) {
	n := len(c.Messages)
	sort.Slice(c.Messages, func(i, j int) bool {
		return c.Messages[i].ReactionsCount > c.Messages[j].ReactionsCount
	})

	playersCount = c.PlayersCount

	var winnersIndexes []int
	winnerReactionsCount := c.Messages[0].ReactionsCount

	index := 0
	for index <= n-1 && c.Messages[index].ReactionsCount == winnerReactionsCount {
		winnersIndexes = append(winnersIndexes, c.Messages[index].ConfigIndex)
		index++
	}

	if len(winnersIndexes) == 1 {
		configIndex = winnersIndexes[0]
		return
	}
	isRandom = true
	randIndex := rand.Intn(len(winnersIndexes))
	configIndex = winnersIndexes[randIndex]
	return
}
