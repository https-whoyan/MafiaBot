package game

import (
	"math/rand"
	"sort"

	botFMT "github.com/https-whoyan/MafiaBot/internal/fmt"
	coreUserPack "github.com/https-whoyan/MafiaBot/internal/user"

	coreGamePack "github.com/https-whoyan/MafiaBot/core/game"
)

var (
	FMTer           = botFMT.DiscordFMTInstance // Same once struct
	ConstRenameMode = coreGamePack.RenameInGuildMode
	VotePing        = 1
)

func GetNewGameConfig(renameProvider *coreUserPack.BotUserRenameProvider) []coreGamePack.GameOption {
	options := []coreGamePack.GameOption{
		coreGamePack.FMTerOpt(FMTer),
		coreGamePack.RenameModeOpt(ConstRenameMode),
		coreGamePack.RenamePrOpt(renameProvider),
		coreGamePack.VotePingOpt(VotePing),
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
	Messages     []*ConfigMessage `redis:"messages"`
}

type ConfigMessage struct {
	ConfigIndex    int    `redis:"configIndex"`
	MessageID      string `json:"messageID"`
	ReactionsCount int
}

func NewConfigMessages(guildID string, peoplesCount int, configsCount int) *ConfigMessages {
	return &ConfigMessages{
		PlayersCount: peoplesCount,
		GuildID:      guildID,
		ConfigsCount: configsCount,
		Messages:     make([]*ConfigMessage, 0),
	}
}

// Setters

func (c *ConfigMessage) SetReactionCount(reactionsCount int) { c.ReactionsCount = reactionsCount }
func (c *ConfigMessages) AddNewMessage(configIndex int, messageID string) {
	c.Messages = append(c.Messages, &ConfigMessage{
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
