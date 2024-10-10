package handlers

import (
	"context"
	"errors"
	"log"
	"strings"

	botChannelPack "github.com/https-whoyan/MafiaBot/internal/channel"
	botFMT "github.com/https-whoyan/MafiaBot/internal/fmt"
	namesPack "github.com/https-whoyan/MafiaBot/internal/handlers/names"
	coreChannelPack "github.com/https-whoyan/MafiaCore/channel"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"
	coreRolePack "github.com/https-whoyan/MafiaCore/roles"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func isCorrectChatID(s *discordgo.Session, chatID string) bool {
	if s.TryLock() {
		defer s.Unlock()
	}

	ch, err := s.Channel(chatID)
	return err == nil && ch != nil
}

// Send message to chatID
func sendMessages(s *discordgo.Session, chatID string, content ...string) (map[string]*discordgo.Message, error) {
	// Represent InteractionMessage by their content.
	messages := make(map[string]*discordgo.Message)
	for _, onceContent := range content {
		message, err := s.ChannelMessageSend(chatID, onceContent)
		if err != nil {
			return nil, err
		}
		messages[onceContent] = message
	}
	return messages, nil
}

// sendToUser Send to userID a message
func sendToUser(s *discordgo.Session, userID string, msg string) error {
	// Create a channelCreate
	channelCreate, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	if channelCreate == nil {
		return errors.New("channelCreate Create Failed, empty channelCreate")
	}
	channelID := channelCreate.ID
	_, err = s.ChannelMessageSend(channelID, msg)
	return err
}

// Response

func Response(s *discordgo.Session, i *discordgo.Interaction, msg string, errLogger *log.Logger) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		errLogger.Println(err)
	}
}

// ____________________
// Error responses
// ____________________

// IsPrivateMessage finds out if a message has been sent to private messages
func IsPrivateMessage(i *discordgo.InteractionCreate) bool {
	return i.GuildID == ""
}

func NoticePrivateChat(s *discordgo.Session, i *discordgo.InteractionCreate,
	fMTer *botFMT.DiscordFMTer, errorLogger *log.Logger) {
	content := fMTer.Bold("All commands are used on the server.") + fMTer.NL() +
		"If you have difficulties in using the bot, " +
		"please refer to the repository documentation: https://github.com/https-whoyan/MafiaBot"
	Response(s, i.Interaction, content, errorLogger)
}

// NoticeIsEmptyGame If game not exists
func NoticeIsEmptyGame(s *discordgo.Session, i *discordgo.InteractionCreate,
	fMTer *botFMT.DiscordFMTer, errorLogger *log.Logger) {
	content := "You can't interact with the game because you haven't registered it" + fMTer.NL() +
		fMTer.Bold("Write the /"+fMTer.Underline(namesPack.RegisterGameCommandName)+" command") + " to start the game."
	Response(s, i.Interaction, content, errorLogger)
}

// __________________
// Channels
// ___________________

// SetRolesChannels to game.
func (c registerGameCommand) setRolesChannels(
	ctx context.Context,
	guildID string, g *coreGamePack.Game,
	s *discordgo.Session,
	f *botFMT.DiscordFMTer,
) ([]string, error) {
	if len(g.GetRoleChannels()) == len(coreRolePack.GetAllNightInteractionRolesNames()) {
		if g.GetMainChannel() != nil {
			// Set it before
			return []string{}, nil
		}
	}
	// Get night interaction roles names
	allRolesNames := coreRolePack.GetAllNightInteractionRolesNames()
	// emptyRolesMp: save not contains channel roles
	emptyRolesMp := make(map[string]bool)
	// mappedRoles: save contains channels roles
	mappedRoles := make(map[string]*botChannelPack.BotRoleChannel)

	addNewChannelIID := func(roleName, channelName string) {
		channelIID, err := c.db.Storage.GetChannelIIDByRole(ctx, guildID, channelName)
		if channelIID == "" {
			emptyRolesMp[channelName] = true
			return
		}
		existsInGuild := checkExistsChannelInDiscord(s, f, channelIID)
		if !existsInGuild {
			emptyRolesMp[channelName] = true
			return
		}
		newRoleChannel, err := botChannelPack.NewBotRoleChannel(c.s, channelIID, roleName)
		if err != nil {
			emptyRolesMp[channelName] = true
			return
		}
		mappedRoles[roleName] = newRoleChannel
	}

	for _, roleName := range allRolesNames {
		if strings.ToLower(roleName) == strings.ToLower(coreRolePack.Don.Name) {
			addNewChannelIID(roleName, coreRolePack.Mafia.Name)
			continue
		}
		addNewChannelIID(roleName, roleName)
	}
	// If a have all roles
	if len(emptyRolesMp) == 0 {
		// Convert
		sliceMappedRoles := lo.Values(mappedRoles)
		interfaceRoleChannelSlice := lo.Map(
			sliceMappedRoles,
			func(ch *botChannelPack.BotRoleChannel, _ int) coreChannelPack.RoleChannel {
				return coreChannelPack.RoleChannel(ch)
			})

		// Save it to g.RoleChannels.
		err := g.SetRoleChannels(interfaceRoleChannelSlice...)

		return []string{}, err
	}

	return lo.Keys(emptyRolesMp), nil
}

// Check, if main channel exists or not
func (c registerGameCommand) existsMainChannel(ctx context.Context, guildID string) bool {
	channelIID, err := c.db.Storage.GetMainChannelIID(ctx, guildID)
	if err != nil {
		return false
	}
	return channelIID != ""
}

func (c registerGameCommand) setMainChannel(ctx context.Context, guildID string, g *coreGamePack.Game) {
	channelIID, _ := c.db.Storage.GetMainChannelIID(ctx, guildID)
	mainChannel, _ := botChannelPack.NewBotMainChannel(c.s, channelIID)
	_ = g.SetMainChannel(mainChannel)
}

func getInfoAboutVote(g *coreGamePack.Game, f *botFMT.DiscordFMTer, vote1 string, vote2 *string) string {
	if vote1 == coreGamePack.EmptyVoteStr {
		return f.IU("You chose not to vote for anyone") + "... 🙄"
	}

	var (
		message string
		players = g.GetActivePlayers()
	)
	if vote2 != nil {
		nick1 := players.SearchPlayerByID(vote1, false).GetNick()
		nick2 := players.SearchPlayerByID(*vote2, false).GetNick()
		message += f.B("You chose to vote for players ")
		message += f.Bl(vote1) + f.B("  :") + f.Bl(nick1) + f.B(" and ") +
			f.Bl(*vote2) + f.B(" :") + f.Bl(nick2)
	} else {
		nick1 := players.SearchPlayerByID(vote1, false).GetNick()
		message += f.B("You chose to vote for player ")
		message += f.Bl(vote1) + f.B("  :") + f.Bl(nick1)
	}
	return message
}

func checkExistsChannelInDiscord(s *discordgo.Session, f *botFMT.DiscordFMTer, channelID string) bool {
	messages, err := sendMessages(
		s, channelID,
		"Bot exists channel test message..."+f.NL()+f.I("Ignore it...."),
	)
	if err != nil {
		return false
	}
	var messageID string
	for _, message := range messages {
		messageID = message.ID
		break
	}
	_ = s.ChannelMessageDelete(channelID, messageID)
	return true
}

// Vote Command validators

func getVotesCountRequired(g *coreGamePack.Game) int {
	nightVoting := g.GetNightVoting()
	if nightVoting == nil {
		return -1
	}
	if !nightVoting.IsTwoVotes {
		return 2
	}
	return 1
}

func (c basicCmd) voteTypeValidator(i *discordgo.Interaction, g *coreGamePack.Game) bool {
	var (
		invalidUsageOfCommand = c.f.BU("Invalid usage of command!")
	)
	gameState := g.GetState()
	votesCountNeed := getVotesCountRequired(g)
	if gameState == coreGamePack.DayState {
		if c.GetName() != namesPack.DayVoteGameCommandName {
			c.response(i, invalidUsageOfCommand)
			return false
		}
		return true
	}
	if votesCountNeed == -1 {
		c.response(i, invalidUsageOfCommand)
		return false
	}
	switch c.GetName() {
	case namesPack.VoteGameCommandName:
		if votesCountNeed == 2 {
			c.response(i, invalidUsageOfCommand)
			return false
		}
	case namesPack.TwoVoteGameCommandName:
		if votesCountNeed == 1 {
			c.response(i, invalidUsageOfCommand)
			return false
		}
	}
	return true
}
