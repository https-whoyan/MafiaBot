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

	coreConfigPack "github.com/https-whoyan/MafiaBot/core/config"
	coreGamePack "github.com/https-whoyan/MafiaBot/core/game"
	coreMessagePack "github.com/https-whoyan/MafiaBot/core/message"
	coreRolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	botCnvPack "github.com/https-whoyan/MafiaBot/internal/converter"
	botFMTPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	botGameCfgPack "github.com/https-whoyan/MafiaBot/internal/game"
	botMsgPack "github.com/https-whoyan/MafiaBot/internal/message"
	botTimeConstsPack "github.com/https-whoyan/MafiaBot/internal/time"

	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
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

func (c AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	if len(i.ApplicationCommandData().Options) == 0 {
		content := "Must be option!"
		Response(s, i, content)
		return
	}

	// Get variables by options
	roleName := i.ApplicationCommandData().Options[0].Name
	requestedChatID := i.ApplicationCommandData().Options[0].Value.(string)

	currDB, containsDB := mongo.GetCurrMongoDB()
	if !containsDB {
		content := "Internal Server Error!"
		Response(s, i, content)
		return
	}

	isFreeChatID, _ := currDB.IsFreeChannelIID(i.GuildID, requestedChatID)
	if !isFreeChatID {
		content := f.B("Provided channel is already used.") + f.NL() + "Please, provide another, available chatID."
		Response(s, i, content)
		return
	}

	if !isCorrectChatID(s, requestedChatID) {
		content := "Invalid Chat IDType!"
		Response(s, i, content)
		return
	}

	err := currDB.SetRoleChannel(i.GuildID, requestedChatID, roleName)
	if err != nil {
		content := "Internal Server Error!"
		Response(s, i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat is used for %v role.", roleName)
	_, err = sendMessages(s, requestedChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + f.NL() + f.Bold("Reuse command.")
		Response(s, i, content)
		return
	}
	messageContent := fmt.Sprintf("Done, now is %v chat is used for %v role.", requestedChatID, roleName) +
		f.NL() + f.BU("Thanks!")
	Response(s, i, messageContent)
}

func (c AddMainChannelCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	if len(i.ApplicationCommandData().Options) == 0 {
		content := "Must be option!"
		Response(s, i, content)
		return
	}

	// Get variables by options
	newChatID := i.ApplicationCommandData().Options[0].Value.(string)

	currDB, containsDB := mongo.GetCurrMongoDB()
	if !containsDB {
		content := "Internal Server Error! Main Chat not add"
		Response(s, i, content)
		return
	}

	isFreeChatID, _ := currDB.IsFreeChannelIID(i.GuildID, newChatID)
	if !isFreeChatID {
		content := f.B("Provided channel is already used.") + f.NL() + "Please, provide another, available chatID."
		Response(s, i, content)
		return
	}

	if !isCorrectChatID(s, newChatID) {
		content := "Invalid Chat IDType!"
		Response(s, i, content)
		return
	}

	err := currDB.SetMainChannel(i.GuildID, newChatID)
	if err != nil {
		content := "Internal Server Error! Main Chat not add"
		Response(s, i, content)
		return
	}

	noticeChatContent := fmt.Sprintf("Now chat will be used as main for games.")
	_, err = sendMessages(s, newChatID, noticeChatContent)
	if err != nil {
		content := "Bot can't write to given chat. " + f.NL() + f.B("Reuse command.")
		Response(s, i, content)
		return
	}
	messageContent := "Done, you configure main interaction chat."
	Response(s, i, messageContent)
}

// _______________
// Game
// _______________

func (c RegisterGameCommand) trySetChannels(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) (correct bool) {
	emptyRoles, err := setRolesChannels(s, i.GuildID, g)
	if err != nil {
		content := "Internal Server Error"
		Response(s, i, content)
		return
	}

	var content string
	if len(emptyRoles) != 0 {
		content += "You don't configure all channels for bot interaction. " +
			"Please, use " + f.BU(AddChannelRoleCommandName) + " to fix " + strings.Join(emptyRoles, ", ") +
			" roles."
	}

	if !existsMainChannel(i.GuildID) {
		if len(content) != 0 {
			content += f.NL() + f.NL()
		}

		content += f.B("You don't configure main channel for bot interaction.") + f.NL() +
			f.I("All messages regarding the game will be sent there.") + f.NL() +
			"To add a main chat channel for interaction, use the command " +
			f.BU(AddMainChannelCommandName)
	}
	setMainChannel(s, i.GuildID, g)
	if len(content) == 0 {
		return true
	}
	Response(s, i, content)
	return false
}

func (c RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	// Validation
	if !c.trySetChannels(s, i, g, f) {

		return
	}

	// Send message.
	Response(s, i, "Ok. InteractionMessage below.")

	// Send additional message and save it IDType
	deadlineStr := strconv.Itoa(botTimeConstsPack.RegistrationDeadlineMinutes)
	responseMessageText := "Registration has begun." + f.NL() +
		f.B("Post "+botFMTPack.RegistrationPlayerSticker+" reaction below.") + f.I(" If you want to be a spectator, "+
		"put the reaction "+botFMTPack.RegistrationSpectatorSticker+".") + f.NL() + f.NL() + f.B(
		f.U("Deadline: "+deadlineStr+" minutes"))

	channelID := i.ChannelID
	message, err := s.ChannelMessageSend(channelID, responseMessageText)
	if err != nil {
		log.Print(err)
	}
	g.SwitchState()
	messageID := message.ID
	// Safe messageID to redis
	currDB, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Print("DB isn't initialed, couldn't get currDB")
		return
	}
	err = currDB.SetInitialGameMessageID(
		i.GuildID, messageID,
		botTimeConstsPack.RegistrationDeadlineSeconds*time.Second)
	if err != nil {
		log.Print(err)
	}
}

func (c FinishGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, _ *botFMTPack.DiscordFMTer) {
	Response(s, i, "Ok... Bad idea, but ok.")
	ch := make(chan coreGamePack.Signal)
	go g.FinishAnyway(ch)
	for signal := range ch {
		log.Println(signal)
	}
}

func (c ChoiceGameConfigCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	currRedisDB, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Println("redis is not exists, command: startGameCommand")
		content := "Internal Server Error!"
		Response(s, i, content)
	}

	registrationMessageID, err := currRedisDB.GetInitialGameMessageID(i.GuildID)
	if (err != nil || registrationMessageID == "") && g.State == coreGamePack.RegisterState {
		messageContent := f.U("Registration Deadline passed!") + f.NL() + "Please, " +
			f.B("use the /"+RegisterGameCommandName+" command") + " to register a new game."
		g.SetState(coreGamePack.NonDefinedState)
		Response(s, i, messageContent)
		return
	}

	// If playersCount not in range [minAvailableCount, maxAvailableCount],
	// Send message that it's impossible to choice config.
	registrationStickerUnicode := botFMTPack.RegistrationPlayerSticker
	_, playersCount := botMsgPack.GetUsersByEmojiID(s, i.ChannelID, registrationMessageID, registrationStickerUnicode)

	allConfigs, nearest, err := coreConfigPack.GetConfigsByPlayersCount(playersCount)
	switch {
	case errors.Is(err, coreConfigPack.SmallCountOfPeopleToConfig):
		content := f.B("The number of players is too small to start the game.") + f.NL() +
			"Number of registered players: " + f.Bl(strconv.Itoa(playersCount)) +
			f.NL() + "Minimum number to vote on game config choices: " + f.Bl(strconv.Itoa(nearest))
		Response(s, i, content)
		return
	case errors.Is(err, coreConfigPack.BigCountOfPeopleToConfig):
		content := f.B("The number of players is too large to start the game.") + f.NL() +
			"Number of registered players: " + f.Bl(strconv.Itoa(playersCount)) + f.NL() +
			"Maximum number to vote on game config choices: " + f.Bl(strconv.Itoa(nearest))
		Response(s, i, content)
		return
	}

	// If playersCount is ok,
	// set empty players to game (to safe it.)
	startBotPlayers, _ := botMsgPack.GetUsersByEmojiID(
		s, i.ChannelID, registrationMessageID, botFMTPack.RegistrationPlayerSticker)
	startGamePlayers := botCnvPack.DiscordUsersToEmptyPlayers(s, i.GuildID,
		startBotPlayers, false)
	g.SetStartPlayers(startGamePlayers)

	// And spectators.
	registerSpectatorStickerUnicode := botFMTPack.RegistrationSpectatorSticker
	botSpectators, _ := botMsgPack.GetUsersByEmojiID(s, i.ChannelID,
		registrationMessageID, registerSpectatorStickerUnicode)
	// We need no duplicates in active and spectators.
	// Then, I get unique Spectators, which not include in active players.
	uniqueBotSpectators, _ := botCnvPack.SetDiff(botSpectators, startBotPlayers)
	spectators := botCnvPack.DiscordUsersToEmptyPlayers(s, i.GuildID, uniqueBotSpectators, true)
	g.SetSpectators(spectators)

	content := "Below is a list of available game configurations. " + f.NL() +
		f.B("If you like the configuration") + ",  please " +
		f.BU("put a reaction ") + botFMTPack.ConfigChoiceSticker + f.BU(" on the post.") + f.NL() +
		"You can give a reaction to any number of configurations. " + f.NL() + f.NL()
	content += f.I("The bot will choose the configuration with the most reactions. "+
		"If there are several configurations, a random one will be selected.") + f.NL() + f.NL()
	deadlineStr := strconv.Itoa(botTimeConstsPack.VotingGameConfigDeadlineMinutes)
	content += "Deadline: " + f.BU(deadlineStr+" minutes") + "."
	Response(s, i, content)

	db, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Println("redis is not exists, command: startGameCommand")
		content = "Internal Server Error!"
		_, _ = sendMessages(s, i.ChannelID, content)
		return
	}

	// Create a ConfigMessages structure to safe all message to redis
	cfgMessages := botGameCfgPack.NewConfigMessages(i.GuildID, playersCount, len(allConfigs))
	for index, config := range allConfigs {
		_, _ = sendMessages(s, i.ChannelID, f.InfoSplitter())
		messageContent := config.GetMessageAboutConfig(f)
		mpMessages, _ := sendMessages(s, i.ChannelID, messageContent)
		cfgMessages.AddNewMessage(index, mpMessages[messageContent].ID)
	}
	err = db.SetConfigGameVotingMessages(cfgMessages, botTimeConstsPack.VotingGameConfigDeadlineSeconds*time.Second)
	if err != nil {
		log.Println(err)
	}

	return
}

func (c StartGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	r, isContains := redis.GetCurrRedisDB()
	if !isContains {
		Response(s, i, "Internal Server Error!")
		return
	}

	cfgMessages, err := r.GetConfigGameVotingMessageID(i.GuildID)
	if err != nil {
		Response(s, i, "Internal Server Error!")
		return
	}

	g.SwitchState()

	// Set Reactions count to cfgMessages
	for _, message := range cfgMessages.Messages {
		messageIID := message.MessageID
		emoji := botFMTPack.ConfigChoiceSticker
		usersWhoLiked, _ := botMsgPack.GetUsersByEmojiID(s, i.ChannelID, messageIID, emoji)

		// Validate, we only want users who participate in the game.
		_, playersMessageCount, err := botCnvPack.GetElsOnlyIncludeFunc(
			usersWhoLiked,
			g.StartPlayers.GetTags(),
			func(u *discordgo.User) string { return u.ID })

		if err != nil {
			Response(s, i, "Internal Server Error!")
			return
		}
		// Set It.
		message.SetReactionCount(playersMessageCount)
	}

	winner, playerCount, isRandom := cfgMessages.GetWinner()
	if isRandom {
		Response(s, i, f.B("Choosing Random configuration..."))
	} else {
		Response(s, i, f.B("Setting the game's configuration..."))
	}
	winnerConfig := coreConfigPack.GetConfigByPlayersCountAndIndex(playerCount, winner)
	err = g.Init(winnerConfig)
	if err != nil {
		_, _ = sendMessages(s, i.ChannelID, f.IU("Can't start game, internal server error!"))
		log.Println(err)
		g.SetState(coreGamePack.RegisterState)
		return
	}
	for _, player := range *g.Active {
		err = SendToUser(s, player.Tag, coreMessagePack.GetStartPlayerDefinition(player, f))
		if err != nil {
			log.Println(err)
		}
	}
	log.Printf("Init Game in %v Guild", i.GuildID)
	gameCh := g.Run(context.Background())
	ProcessGameChan(g, f, gameCh)
}

// _______________
// Voting
// _______________

func (c GameVoteCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	options := i.ApplicationCommandData().Options
	vote := options[0].Value.(string)
	vP := coreGamePack.NewVoteProvider(i.Member.User.ID, vote, true)
	err := g.NightVoteValidator(vP, nil)

	switch {
	case err == nil:
		goto answerToVote
	case errors.Is(err, coreGamePack.IncorrectVotedPlayer):
		message := f.B("Right now you are not allowed to leave your vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTime):
		message := f.B("It is not possible to apply a vote not during the night or not during daytime voting.") +
			f.NL() + f.U("Use the command later.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteType):
		message := f.B("Incorrect format for entering the ID of the player you are voting for.") + f.NL() +
			"Available options " + f.I("live players") + ":" + f.NL() + f.Tab()
		var allIDS []string
		for _, player := range *g.Active {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.PlayerIsMutedErr):
		message := f.B("Oops, you are muted now.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerIsNotAlive):
		message := f.B("I think the dead can't vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerNotFound):
		message := f.B(fmt.Sprintf("Player IDType %v is not found alive.", vote)) + f.NL()
		message += f.NL() + "Available options " + f.I("live players") + ":" + f.NL() + f.Tab()
		var allIDS []string
		for _, player := range *g.Active {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePingError):
		message := f.BU("You've already voted for this player recently.") + f.NL() + f.Tab()
		message += f.I(fmt.Sprintf("You cannot vote for the same player 2 times within %v nights.",
			g.VotePing+1)) + f.NL() + f.NL()
		message += f.IU("Please, re-vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.TwoVoteRequiredErr):
		message := f.B("Your role requires a double vote.") + f.NL()
		message += f.BU("Use /" + TwoVoteGameCommandName + " command.")
		Response(s, i, message)
		return
	}
answerToVote:
	Response(s, i, f.B("Your vote counts. "))
	g.VoteChan <- vP

}

func (c GameTwoVoteCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	options := i.ApplicationCommandData().Options
	vote1 := options[0].Value.(string)
	vote2 := options[1].Value.(string)
	vP := coreGamePack.NewTwoVoteProvider(i.Member.User.ID, vote1, vote2, true)
	err := g.NightTwoVoteValidator(vP, nil)

	switch {
	case err == nil:
		goto answerToTwoVote
	case errors.Is(err, coreGamePack.IncorrectVotedPlayer):
		message := f.B("Right now you are not allowed to leave your vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTime):
		message := f.B("It is not possible to apply a vote not during the night or not during daytime voting.") +
			f.NL() + f.U("Use the command later.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.OneVoteRequiredErr):
		message := "Your role requires one vote, not 2." + f.NL()
		message += "Use " + f.B("/"+VoteGameCommandName) + " command"
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteType):
		message := f.B("Incorrect format for entering the IDType of the player you are voting for.") + f.NL() +
			"Available options " + f.I("live players") + ":" + f.NL() + f.Tab()
		var allIDS []string
		for _, player := range *g.Active {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.PlayerIsMutedErr):
		message := f.B("Oops, you are muted now.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerIsNotAlive):
		message := f.B("I think the dead  can't vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerNotFound):
		message := f.B(fmt.Sprintf("Player IDType %v or %v is not found alive.", vote1, vote2)) + f.NL()
		message += f.NL() + "Available options " + f.I("live players") + ":" + f.NL() + f.Tab()
		var allIDS []string
		for _, player := range *g.Active {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePingError):
		message := f.BU("You've already voted for this player recently.") + f.NL() + f.Tab()
		message += f.I(fmt.Sprintf("You cannot vote for the same player 2 times within %v nights.",
			g.VotePing+1)) + f.NL() + f.NL()
		message += f.IU("Please, re-vote.")
		Response(s, i, message)
		return
	}
answerToTwoVote:
	Response(s, i, f.B("Your vote counts. "))
	g.TwoVoteChan <- vP
}

func (c DayVoteCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	g *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	options := i.ApplicationCommandData().Options
	vote := options[0].Value.(string)
	fmt.Println(518)
	vP := coreGamePack.NewVoteProvider(i.Member.User.ID, vote, true)
	fmt.Println(g.Active)
	err := g.DayVoteValidator(vP)
	fmt.Println("ti dalbaeb, ", err)

	switch {
	case err == nil:
		goto answerToVote
	case errors.Is(err, coreGamePack.IncorrectVotedPlayer):
		message := f.B("Right now you are not allowed to leave your vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteTime):
		message := f.B("It is not possible to apply a vote not during the night or not during daytime voting.") +
			f.NL() + f.U("Use the command later.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.IncorrectVoteType):
		message := f.B("Incorrect format for entering the ID of the p you are voting for.") + f.NL() +
			"Available options " + f.I("live players") + ":" + f.NL() + f.Tab()
		var allIDS []string
		for _, p := range *g.Active {
			allIDS = append(allIDS, fmt.Sprintf("%v", p.ID))
		}
		message += strings.Join(allIDS, ", ")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.PlayerIsMutedErr):
		message := f.B("Oops, you are muted now.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerIsNotAlive):
		message := f.B("I think the dead can't vote.")
		Response(s, i, message)
		return
	case errors.Is(err, coreGamePack.VotePlayerNotFound):
		message := f.B(fmt.Sprintf("Player IDType %v is not found alive.", vote)) + f.NL()
		message += f.NL() + "Available options " + f.I("live players") + ":" + f.NL() + f.Tab()
		var allIDS []string
		for _, player := range *g.Active {
			allIDS = append(allIDS, fmt.Sprintf("%v", player.ID))
		}
		message += strings.Join(allIDS, ", ")
		Response(s, i, message)
		return
	}
answerToVote:
	Response(s, i, f.B("Your day vote been accepted."))
	g.VoteChan <- vP
}

// _____________
// Other
// _____________

func (c YanLohCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, _ *botFMTPack.DiscordFMTer) {
	messageContent := "Не, лол, что бы его же ботом его обзывать..."
	Response(s, i, messageContent)

	// async kick requester
	guildID := i.GuildID
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(sessId, kickedUserID string) {
		defer wg.Done()
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err := s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Member.User.ID)
	wg.Wait()
}

func (c AboutRolesCommand) Execute(s *discordgo.Session, i *discordgo.Interaction,
	_ *coreGamePack.Game, f *botFMTPack.DiscordFMTer) {
	messageContent := f.Bold("Below information about all roles:") + f.NL()
	Response(s, i, messageContent)

	sendMessage := func(s *discordgo.Session, i *discordgo.Interaction, message string) {
		_, err := s.ChannelMessageSend(i.ChannelID, message)
		if err != nil {
			log.Print(err)
		}
	}

	messages := coreRolesPack.GetDefinitionsOfAllRoles(f, 2000)
	for _, message := range messages {
		sendMessage(s, i, coreRolesPack.FixDescription(message))
	}
}
