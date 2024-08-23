package handlers

import (
	"errors"
	"github.com/https-whoyan/MafiaCore/channel"
	"github.com/samber/lo"
	"log"
	"strings"

	botChannelPack "github.com/https-whoyan/MafiaBot/internal/channel"
	botFMT "github.com/https-whoyan/MafiaBot/internal/fmt"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"
	coreRolePack "github.com/https-whoyan/MafiaCore/roles"

	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"

	"github.com/bwmarrin/discordgo"
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

// SendToUser Send to userID a message
func SendToUser(s *discordgo.Session, userID string, msg string) error {
	// Create a channel
	channel, err := s.UserChannelCreate(userID)
	if err != nil || channel == nil {
		if channel == nil {
			return errors.New("channel Create Failed, empty channel")
		}
		return err
	}
	channelID := channel.ID
	_, err = s.ChannelMessageSend(channelID, msg)
	return err
}

// ___________________
// Response func
// ___________________

// Response reply to interaction by provided content (s.InteractionResponse)
func Response(s *discordgo.Session, i *discordgo.Interaction, content string) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Print(err)
	}
}

// ____________________
// Error responses
// ____________________

// IsPrivateMessage finds out if a message has been sent to private messages
func IsPrivateMessage(i *discordgo.InteractionCreate) bool {
	return i.GuildID == ""
}

func NoticePrivateChat(s *discordgo.Session, i *discordgo.InteractionCreate, fMTer *botFMT.DiscordFMTer) {
	content := fMTer.Bold("All commands are used on the server.") + fMTer.NL() +
		"If you have difficulties in using the bot, " +
		"please refer to the repository documentation: https://github.com/https-whoyan/MafiaBot"
	Response(s, i.Interaction, content)
}

// NoticeIsEmptyGame If game not exists
func NoticeIsEmptyGame(s *discordgo.Session, i *discordgo.InteractionCreate, fMTer *botFMT.DiscordFMTer) {
	content := "You can't interact with the game because you haven't registered it" + fMTer.NL() +
		fMTer.Bold("Write the /"+fMTer.Underline(RegisterGameCommandName)+" command") + " to start the game."
	Response(s, i.Interaction, content)
}

// __________________
// Channels
// ___________________

// SetRolesChannels to game.
func setRolesChannels(s *discordgo.Session, guildID string, g *coreGamePack.Game) ([]string, error) {
	if len(g.GetRoleChannels()) == len(coreRolePack.GetAllNightInteractionRolesNames()) {
		if g.GetMainChannel() != nil {
			// Set it before
			return []string{}, nil
		}
	}
	// Get night interaction roles names
	allRolesNames := coreRolePack.GetAllNightInteractionRolesNames()
	// Get curr MongoDB struct
	currDB, isContains := mongo.GetCurrMongoDB()
	if !isContains {
		return []string{}, errors.New("MongoDB doesn't initialized")
	}
	// emptyRolesMp: save not contains channel roles
	emptyRolesMp := make(map[string]bool)
	// mappedRoles: save contains channels roles
	mappedRoles := make(map[string]*botChannelPack.BotRoleChannel)

	addNewChannelIID := func(roleName, channelName string) {
		channelIID, err := currDB.GetChannelIIDByRole(guildID, channelName)
		if channelIID == "" {
			emptyRolesMp[channelName] = true
			return
		}
		newRoleChannel, err := botChannelPack.NewBotRoleChannel(s, channelIID, roleName)
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
			func(ch *botChannelPack.BotRoleChannel, _ int) channel.RoleChannel {
				return channel.RoleChannel(ch)
			})

		// Save it to g.RoleChannels.
		err := g.SetRoleChannels(interfaceRoleChannelSlice...)

		return []string{}, err
	}

	return lo.Keys(emptyRolesMp), nil
}

// Check, if main channel exists or not
func existsMainChannel(guildID string) bool {
	currMongo, exists := mongo.GetCurrMongoDB()
	if !exists {
		return false
	}
	channelIID, err := currMongo.GetMainChannelIID(guildID)
	if err != nil {
		return false
	}
	return channelIID != ""
}

func setMainChannel(s *discordgo.Session, guildID string, g *coreGamePack.Game) {
	currMongo, _ := mongo.GetCurrMongoDB()
	channelIID, _ := currMongo.GetMainChannelIID(guildID)
	mainChannel, _ := botChannelPack.NewBotMainChannel(s, channelIID)
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

func voteTypeValidator(s *discordgo.Session, i *discordgo.Interaction,
	c Command, f *botFMT.DiscordFMTer, g *coreGamePack.Game) bool {
	var (
		notRequiredVoteNowMessage = "Right now, " + f.B("the game is not in voting mode.") +
			f.NL() + "Please repeat the command later."
		requiredOneVote = f.B("It's a one-vote goal right now. ") + "Please try again later."
		requiredTwoVote = f.B("It's a two-vote goal right now. ") + "Please try again later."
		dayVoteRequired = f.B("It's a day-vote goal right now. ")
	)
	gameState := g.GetState()
	votesCountNeed := getVotesCountRequired(g)
	if gameState == coreGamePack.DayState {
		if c.GetName() != DayVoteGameCommandName {
			Response(s, i, dayVoteRequired)
			return false
		}
		return true
	}
	if votesCountNeed == -1 {
		Response(s, i, notRequiredVoteNowMessage)
		return false
	}
	switch c.GetName() {
	case VoteGameCommandName:
		if votesCountNeed == 2 {
			Response(s, i, requiredOneVote)
			return false
		}
	case TwoVoteGameCommandName:
		if votesCountNeed == 1 {
			Response(s, i, requiredTwoVote)
			return false
		}
	}
	return true
}
