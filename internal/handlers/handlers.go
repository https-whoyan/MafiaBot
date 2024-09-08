package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	botCnvPack "github.com/https-whoyan/MafiaBot/internal/converter"
	botFMTPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	botGameCfgPack "github.com/https-whoyan/MafiaBot/internal/game"
	botMsgPack "github.com/https-whoyan/MafiaBot/internal/message"
	botTimeConstsPack "github.com/https-whoyan/MafiaBot/internal/time"
	coreConfigPack "github.com/https-whoyan/MafiaCore/config"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"
	coreMessagePack "github.com/https-whoyan/MafiaCore/message"
	coreRolesPack "github.com/https-whoyan/MafiaCore/roles"

	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"

	"github.com/bwmarrin/discordgo"
)

// _________________________________
// ********************************
// This contains all bot handlers.
//********************************
// _________________________________

// __________________
// Channels
// __________________

const (
	InternalServerErrorMessage = "Internal Server error!"
)

func (c addChannelRoleCommand) Execute(ctx context.Context, i *discordgo.Interaction, _ *coreGamePack.Game) {
	if len(i.ApplicationCommandData().Options) == 0 {
		content := "Must be option!"
		c.response(i, content)
		return
	}

	// Get variables by options
	roleName := i.ApplicationCommandData().Options[0].Name
	requestedChatID := i.ApplicationCommandData().Options[0].Value.(string)

	isFreeChatID, _ := c.storage.IsFreeChannelIID(ctx, i.GuildID, requestedChatID)
	if !isFreeChatID {
		content := c.f.B("Provided channel is already used.") + c.f.NL() + "Please, provide another, available chatID."
		c.response(i, content)
		return
	}

	if !isCorrectChatID(c.s, requestedChatID) {
		content := "Invalid Chat ID!"
		c.response(i, content)
		return
	}

	err := c.storage.SetRoleChannel(ctx, i.GuildID, requestedChatID, roleName)
	if err != nil {
		content := "Internal Server Error!"
		c.response(i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat is used for %v role.", roleName)
	_, err = sendMessages(c.s, requestedChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + c.f.NL() + c.f.Bold("Reuse command.")
		c.response(i, content)
		return
	}
	messageContent := fmt.Sprintf("Done, now is %v chat is used for %v role.", requestedChatID, roleName) +
		c.f.NL() + c.f.BU("Thanks!")
	c.response(i, messageContent)
}

func (c addMainChannelCommand) Execute(ctx context.Context, i *discordgo.Interaction, _ *coreGamePack.Game) {
	if len(i.ApplicationCommandData().Options) == 0 {
		content := "Must be option!"
		c.response(i, content)
		return
	}

	// Get variables by options
	newChatID := i.ApplicationCommandData().Options[0].Value.(string)
	isFreeChatID, _ := c.storage.IsFreeChannelIID(ctx, i.GuildID, newChatID)
	if !isFreeChatID {
		content := c.f.B("Provided channel is already used.") + c.f.NL() + "Please, provide another, available chatID."
		c.response(i, content)
		return
	}

	if !isCorrectChatID(c.s, newChatID) {
		content := "Invalid Chat ID!"
		c.response(i, content)
		return
	}

	err := c.storage.SetMainChannel(ctx, i.GuildID, newChatID)
	if err != nil {
		content := "Internal Server Error! Main Chat not add"
		c.response(i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat will be used as main for games.")
	_, err = sendMessages(c.s, newChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + c.f.NL() + c.f.B("Reuse command.")
		c.response(i, content)
		return
	}
	messageContent := "Done, you configure main interaction chat."
	c.response(i, messageContent)
}

// _______________
// Game
// _______________

func (c registerGameCommand) trySetChannels(ctx context.Context, i *discordgo.Interaction, g *coreGamePack.Game) (correct bool) {
	emptyRoles, err := c.setRolesChannels(ctx, i.GuildID, g)
	if err != nil {
		content := "Internal Server Error"
		c.response(i, content)
		return
	}

	var content string
	if len(emptyRoles) != 0 {
		content += "You don't configure all channels for bot interaction. " +
			"Please, use " + c.f.BU(addChannelRoleCommandName) + " to fix " + strings.Join(emptyRoles, ", ") +
			" roles."
	}

	if !c.existsMainChannel(ctx, i.GuildID) {
		if len(content) != 0 {
			content += c.f.NL() + c.f.NL()
		}

		content += c.f.B("You don't configure main channel for bot interaction.") + c.f.NL() +
			c.f.I("All messages regarding the game will be sent there.") + c.f.NL() +
			"To add a main chat channel for interaction, use the command " +
			c.f.BU(addMainChannelCommandName)
	}
	c.setMainChannel(ctx, i.GuildID, g)
	if len(content) == 0 {
		return true
	}
	c.response(i, content)
	return false
}

func (c registerGameCommand) Execute(ctx context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	// Validation
	if !c.trySetChannels(ctx, i, g) {
		return
	}
	// Send message.
	c.response(i, "Ok. InteractionMessage below.")
	// Send additional message and save it ID
	deadlineStr := strconv.Itoa(botTimeConstsPack.RegistrationDeadlineMinutes)
	responseMessageText := "Registration has begun." + c.f.NL() +
		c.f.B("Post "+botFMTPack.RegistrationPlayerSticker+" reaction below.") + c.f.I(" If you want to be a spectator, "+
		"put the reaction "+botFMTPack.RegistrationSpectatorSticker+".") + c.f.NL() + c.f.NL() + c.f.B(
		c.f.U("Deadline: "+deadlineStr+" minutes"),
	)

	channelID := i.ChannelID
	message, _ := c.s.ChannelMessageSend(channelID, responseMessageText)
	g.SwitchState()
	messageID := message.ID
	// Safe messageID to redis
	_ = c.db.Hasher.SetInitialGameMessageID(
		ctx,
		i.GuildID, messageID,
		botTimeConstsPack.RegistrationDeadlineSeconds*time.Second,
	)
	_ = c.db.Hasher.SetChannelStorage(
		ctx,
		i.GuildID, i.ChannelID, redis.SetInitialGameStorage,
		botTimeConstsPack.RegistrationDeadlineSeconds*time.Second,
	)
}

func (c finishGameCommand) Execute(_ context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	c.response(i, "Ok... Bad idea, but ok.")
	g.FinishAnyway()
}

func (c choiceGameConfigCommand) Execute(ctx context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	registrationMessageID, err := c.hasher.GetInitialGameMessageID(ctx, i.GuildID)
	if (err != nil || registrationMessageID == "") && g.GetState() == coreGamePack.RegisterState {
		messageContent := c.f.U("Registration Deadline passed!") + c.f.NL() + "Please, " +
			c.f.B("use the /"+RegisterGameCommandName+" command") + " to register a new game."
		g.SetState(coreGamePack.NonDefinedState)
		c.response(i, messageContent)
		return
	}

	// If playersCount not in range [minAvailableCount, maxAvailableCount],
	// Send message that it's impossible to choice config.
	registrationStickerUnicode := botFMTPack.RegistrationPlayerSticker
	initialMessageChannelIID, err := c.hasher.GetChannelStorage(ctx, i.GuildID, redis.SetInitialGameStorage)
	if err != nil || initialMessageChannelIID == "" {
		content := "Internal Server Error!"
		c.response(i, content)
	}
	_, playersCount := botMsgPack.GetUsersByEmojiID(c.s, initialMessageChannelIID, registrationMessageID, registrationStickerUnicode)

	allConfigs, nearest, err := coreConfigPack.GetConfigsByPlayersCount(playersCount)
	switch {
	case errors.Is(err, coreConfigPack.SmallCountOfPeopleToConfig):
		content := c.f.B("The number of players is too small to start the game.") + c.f.NL() +
			"Number of registered players: " + c.f.Bl(strconv.Itoa(playersCount)) +
			c.f.NL() + "Minimum number to vote on game config choices: " + c.f.Bl(strconv.Itoa(nearest))
		c.response(i, content)
		return
	case errors.Is(err, coreConfigPack.BigCountOfPeopleToConfig):
		content := c.f.B("The number of players is too large to start the game.") + c.f.NL() +
			"Number of registered players: " + c.f.Bl(strconv.Itoa(playersCount)) + c.f.NL() +
			"Maximum number to vote on game config choices: " + c.f.Bl(strconv.Itoa(nearest))
		c.response(i, content)
		return
	}

	// If playersCount is ok,
	// set empty players to game (to safe it.)
	startBotPlayers, _ := botMsgPack.GetUsersByEmojiID(
		c.s, i.ChannelID, registrationMessageID, botFMTPack.RegistrationPlayerSticker)
	startGamePlayers := botCnvPack.DiscordUsersToEmptyPlayers(c.s, i.GuildID,
		startBotPlayers, false)
	g.SetStartPlayers(startGamePlayers)

	// And spectators.
	registerSpectatorStickerUnicode := botFMTPack.RegistrationSpectatorSticker
	botSpectators, _ := botMsgPack.GetUsersByEmojiID(c.s, i.ChannelID,
		registrationMessageID, registerSpectatorStickerUnicode)
	// We need no duplicates in active and spectators.
	// Then, I get unique Spectators, which not include in active players.
	uniqueBotSpectators, _ := botCnvPack.SetDiff(botSpectators, startBotPlayers)
	spectators := botCnvPack.DiscordUsersToEmptyPlayers(c.s, i.GuildID, uniqueBotSpectators, true)
	g.SetSpectators(spectators)

	content := "Below is a list of available game configurations. " + c.f.NL() +
		c.f.B("If you like the configuration") + ",  please " +
		c.f.BU("put a reaction ") + botFMTPack.ConfigChoiceSticker + c.f.BU(" on the post.") + c.f.NL() +
		"You can give a reaction to any number of configurations. " + c.f.NL() + c.f.NL()
	content += c.f.I("The bot will choose the configuration with the most reactions. "+
		"If there are several configurations, a random one will be selected.") + c.f.NL() + c.f.NL()
	deadlineStr := strconv.Itoa(botTimeConstsPack.VotingGameConfigDeadlineMinutes)
	content += "Deadline: " + c.f.BU(deadlineStr+" minutes") + "."
	c.response(i, content)

	// Create a ConfigMessages structure to safe all message to redis
	cfgMessages := botGameCfgPack.NewConfigMessages(i.GuildID, playersCount, len(allConfigs))
	for index, config := range allConfigs {
		_, _ = sendMessages(c.s, i.ChannelID, c.f.InfoSplitter())
		messageContent := config.GetMessageAboutConfig(c.f)
		mpMessages, _ := sendMessages(c.s, i.ChannelID, messageContent)
		cfgMessages.AddNewMessage(index, mpMessages[messageContent].ID)
	}
	err = c.hasher.SetConfigGameVotingMessages(ctx, cfgMessages, botTimeConstsPack.VotingGameConfigDeadlineSeconds*time.Second)
	if err != nil {
		log.Println(err)
	}

	return
}

func (c startGameCommand) Execute(ctx context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	cfgMessages, err := c.hasher.GetConfigGameVotingMessageID(ctx, i.GuildID)
	if err != nil {
		c.response(i, InternalServerErrorMessage)
		return
	}

	g.SwitchState()

	// Set Reactions count to cfgMessages
	for _, message := range cfgMessages.Messages {
		messageIID := message.MessageID
		emoji := botFMTPack.ConfigChoiceSticker
		usersWhoLiked, _ := botMsgPack.GetUsersByEmojiID(c.s, i.ChannelID, messageIID, emoji)

		// Validate, we only want users who participate in the game.
		players := g.GetStartPlayers()
		_, playersMessageCount, internalErr := botCnvPack.GetElsOnlyIncludeFunc(
			usersWhoLiked,
			players.GetTags(),
			func(u *discordgo.User) string { return u.ID })

		if internalErr != nil {
			c.response(i, InternalServerErrorMessage)
			return
		}
		// Set It.
		message.SetReactionCount(playersMessageCount)
	}

	winner, playerCount, isRandom := cfgMessages.GetWinner()
	if isRandom {
		c.response(i, c.f.B("Choosing Random configuration..."))
	} else {
		c.response(i, c.f.B("Setting the game's configuration..."))
	}
	winnerConfig := coreConfigPack.GetConfigByPlayersCountAndIndex(playerCount, winner)
	err = g.Init(winnerConfig)
	if err != nil {
		_, _ = sendMessages(c.s, i.ChannelID, c.f.IU("Can't start game, internal server error!"))
		log.Println(err)
		g.SetState(coreGamePack.RegisterState)
		return
	}
	for _, player := range g.GetActivePlayers() {
		err = SendToUser(c.s, player.Tag, coreMessagePack.GetStartPlayerDefinition(player, c.f))
		if err != nil {
			log.Println(err)
		}
	}
	log.Printf("Init Game in %v Guild", i.GuildID)
	gameCh := g.Run(context.Background())
	ProcessGameChan(g, c.f, gameCh)
}

// _______________
// Voting
// _______________

func (c gameVoteCommand) Execute(_ context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	if ok := c.voteTypeValidator(i, g); !ok {
		return
	}
	options := i.ApplicationCommandData().Options
	vote := options[0].Value.(string)
	vP := coreGamePack.NewVoteProvider(i.Member.User.ID, vote, true, false)
	err := g.SetNightVote(vP)

	switch {
	case err == nil:
		var message = "Your vote counts." + c.f.NL()
		message += getInfoAboutVote(g, c.f, vote, nil)
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTimeErr):
		message := c.f.B("Right now you are not allowed to leave your vote.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTimeErr):
		message := c.f.B("It is not possible to apply a vote not during the night or not during daytime voting.") +
			c.f.NL() + c.f.U("Use the command later.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteType):
		message := c.f.B("Incorrect format for entering the ID of the player you are voting for.") + c.f.NL() +
			"Available options " + c.f.I("live players") + ":" + c.f.NL() + c.f.Tab()
		var allIDS []string
		for _, player := range g.GetActivePlayers() {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.PlayerIsMutedErr):
		message := c.f.B("Oops, you are muted now.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerIsNotAliveErr):
		message := c.f.B("I think the dead can't vote.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerNotFoundErr):
		message := c.f.B(fmt.Sprintf("Player ID %v is not found alive.", vote)) + c.f.NL()
		message += c.f.NL() + "Available options " + c.f.I("live players") + ":" + c.f.NL() + c.f.Tab()
		var allIDS []string
		for _, player := range g.GetActivePlayers() {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePingErr):
		message := c.f.BU("You've already voted for this player recently.") + c.f.NL() + c.f.Tab()
		message += c.f.I(fmt.Sprintf("You cannot vote for the same player 2 times within %v nights.",
			g.GetVotePing()+1)) + c.f.NL() + c.f.NL()
		message += c.f.IU("Please, re-vote.")
		c.response(i, message)
		return
	}
}

func (c gameTwoVoteCommand) Execute(_ context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	if ok := c.voteTypeValidator(i, g); !ok {
		return
	}
	options := i.ApplicationCommandData().Options
	vote1 := options[0].Value.(string)
	vote2 := options[1].Value.(string)
	vP := coreGamePack.NewTwoVoteProvider(i.Member.User.ID, vote1, vote2, true, false)
	err := g.SetNightTwoVote(vP)

	switch {
	case err == nil:
		var message = "Your vote counts." + c.f.NL()
		message += getInfoAboutVote(g, c.f, vote1, &vote2)
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTimeErr):
		message := c.f.B("Right now you are not allowed to leave your vote.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteType):
		message := c.f.B("Incorrect format for entering the ID of the player you are voting for.") + c.f.NL() +
			"Available options " + c.f.I("live players") + ":" + c.f.NL() + c.f.Tab()
		var allIDS []string
		for _, player := range g.GetActivePlayers() {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.PlayerIsMutedErr):
		message := c.f.B("Oops, you are muted now.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerIsNotAliveErr):
		message := c.f.B("I think the dead player can't vote...")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerNotFoundErr):
		message := c.f.B(fmt.Sprintf("Player ID %v or %v is not found alive.", vote1, vote2)) + c.f.NL()
		message += c.f.NL() + "Available options " + c.f.I("live players") + ":" + c.f.NL() + c.f.Tab()
		var allIDS []string
		for _, player := range g.GetActivePlayers() {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePingErr):
		message := c.f.BU("You've already voted for this player recently.") + c.f.NL() + c.f.Tab()
		message += c.f.I(fmt.Sprintf("You cannot vote for the same player 2 times within %v nights.",
			g.GetVotePing()+1)) + c.f.NL() + c.f.NL()
		message += c.f.IU("Please, re-vote.")
		c.response(i, message)
		return
	}
}

func (c dayVoteCommand) Execute(_ context.Context, i *discordgo.Interaction, g *coreGamePack.Game) {
	g.SwitchState()
	if ok := c.voteTypeValidator(i, g); !ok {
		return
	}
	options := i.ApplicationCommandData().Options
	vote := options[0].Value.(string)
	vP := coreGamePack.NewVoteProvider(i.Member.User.ID, vote, true, false)
	err := g.SetDayVote(vP)

	switch {
	case err == nil:
		var message = "Your vote counts." + c.f.NL()
		message += getInfoAboutVote(g, c.f, vote, nil)
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTimeErr):
		message := c.f.B("It is not possible to apply a vote not during the night or not during daytime voting.") +
			c.f.NL() + c.f.U("Use the command later.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteType):
		message := c.f.B("Incorrect format for entering the ID of the p you are voting for.") + c.f.NL() +
			"Available options " + c.f.I("live players") + ":" + c.f.NL() + c.f.Tab()
		var allIDS []string
		for _, p := range g.GetActivePlayers() {
			allIDS = append(allIDS, fmt.Sprintf("%v", p.ID))
		}
		message += strings.Join(allIDS, ", ")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.PlayerIsMutedErr):
		message := c.f.B("Oops, you are muted now.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerIsNotAliveErr):
		message := c.f.B("I think the dead can't vote.")
		c.response(i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerNotFoundErr):
		message := c.f.B(fmt.Sprintf("Player ID %v is not found alive.", vote)) + c.f.NL()
		message += c.f.NL() + "Available options " + c.f.I("live players") + ":" + c.f.NL() + c.f.Tab()
		var allIDS []string
		for _, player := range g.GetActivePlayers() {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		c.response(i, message)
		return
	}
}

// _____________
// Other
// _____________

func (c yanLohCommand) Execute(_ context.Context, i *discordgo.Interaction, _ *coreGamePack.Game) {
	messageContent := "Не, лол, что бы его же ботом его обзывать..."
	c.response(i, messageContent)

	// async kick requester
	guildID := i.GuildID
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(sessId, kickedUserID string) {
		defer wg.Done()
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err := c.s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Member.User.ID)
	wg.Wait()
}

func (c aboutRolesCommand) Execute(context context.Context, i *discordgo.Interaction, _ *coreGamePack.Game) {
	messageContent := c.f.Bold("Below information about all roles:") + c.f.NL()
	c.response(i, messageContent)

	sendMessage := func(s *discordgo.Session, i *discordgo.Interaction, message string) {
		_, err := s.ChannelMessageSend(i.ChannelID, message)
		if err != nil {
			log.Print(err)
		}
	}

	messages := coreRolesPack.GetDefinitionsOfAllRoles(c.f, 2000)
	for _, message := range messages {
		sendMessage(c.s, i, coreRolesPack.FixDescription(message))
	}
}
