package bot

import (
	"context"
	"log"
	"os"
	"sync"

	botFMTPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	botGamePack "github.com/https-whoyan/MafiaBot/internal/game"
	handlerPack "github.com/https-whoyan/MafiaBot/internal/handlers"
	userPack "github.com/https-whoyan/MafiaBot/internal/user"
	pkgPack "github.com/https-whoyan/MafiaBot/pkg"
	gamePack "github.com/https-whoyan/MafiaCore/game"

	"github.com/bwmarrin/discordgo"
)

// ____________
// BotConfig
// ____________

type Config struct {
	token string
}

func LoadBotConfig() *Config {
	token := os.Getenv("BOT_TOKEN")
	return &Config{
		token: token,
	}
}

// ________
// Bot
// ________

type Bot struct {
	sync.Mutex
	// DiscordGo token
	token string
	// DiscordGo Session
	Session *discordgo.Session
	// Seen from https://github.com/bwmarrin/discordgo/tree/master/examples/slash_commands
	// The key is the name of the command.
	Commands map[string]handlerPack.Command
	// Games this a map,
	// the key in which is the State of the server where the bot is running,
	// and the value is the game.
	Games map[string]*gamePack.Game
	// To save DiscordGo.ApplicationCommand's for closing deleting.
	registeredCommands []*discordgo.ApplicationCommand
	// To format messages.
	//
	// Implement of FmtInterface.
	FMTer *botFMTPack.DiscordFMTer
	// Databases
	Databases *pkgPack.Database
}

func InitBot(ctx context.Context, cnf *Config, databases *pkgPack.Database) (*Bot, error) {
	token := cnf.token
	botStr := "Bot " + token
	s, err := discordgo.New(botStr)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		token:     token,
		Session:   s,
		Commands:  make(map[string]handlerPack.Command),
		Games:     make(map[string]*gamePack.Game),
		FMTer:     botFMTPack.DiscordFMTInstance,
		Databases: databases,
	}

	bot.initBotCommands()
	bot.registerHandlers(ctx)
	return bot, nil
}

func (b *Bot) loginAs() {
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
}

func (b *Bot) Open() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}
	b.loginAs()
	b.registerCommands()
	return nil
}

func (b *Bot) Close() error {
	b.finishAllGames()
	b.deleteHandlers()
	b.removeRegisteredCommands()
	err := b.Session.Close()
	if err != nil {
		return err
	}
	return nil
}

// ____________________________________________________________
// All below functions initialize or delete the required variables.
// ____________________________________________________________

func (b *Bot) initCommand(c handlerPack.Command) {
	commandName := c.GetName()
	b.Commands[commandName] = c
}

func (b *Bot) initBotCommands() {
	// Channels
	b.initCommand(handlerPack.NewAddMainChannelCommand(b.Session, b.Databases.Storage))
	b.initCommand(handlerPack.NewAddChannelRoleCommand(b.Session, b.Databases.Storage))

	// Game
	b.initCommand(handlerPack.NewRegisterGameCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewChoiceGameConfigCommand(b.Session, b.Databases.Hasher))
	b.initCommand(handlerPack.NewStartGameCommand(b.Session, b.Databases.Hasher))

	// Vote
	b.initCommand(handlerPack.NewGameVoteCommand(b.Session))
	b.initCommand(handlerPack.NewGameTwoVoteCommand(b.Session))
	b.initCommand(handlerPack.NewDayVoteCommand(b.Session))

	// finish game
	b.initCommand(handlerPack.NewFinishGameCommand(b.Session))

	// Other
	b.initCommand(handlerPack.NewYanLohCommand(b.Session))
	b.initCommand(handlerPack.NewAboutRolesCommand(b.Session))

}

func (b *Bot) registerHandlers(ctx context.Context) {
	log.Print("Register handlers")
	for _, cmd := range b.Commands {
		// To avoid closing the loop
		newCmd := cmd
		cmdName := newCmd.GetName()
		log.Printf("Register handler, command name: %v", cmdName)
		// Lock
		newHandler := b.getSIHandler(ctx, newCmd, cmdName)
		b.Session.AddHandler(newHandler)
	}
}

// getSIHandler Is bot command handler get Function.
// All comments in function.
func (b *Bot) getSIHandler(ctx context.Context, cmd handlerPack.Command, cmdName string) func(
	s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)

				return
			}
		}()
		// I recognize the name of the team
		executedCommandName := i.ApplicationCommandData().Name

		// If it not equals as iterable cmd.Name
		if executedCommandName != cmdName {
			return
		}

		// If it executed in private chat
		if handlerPack.IsPrivateMessage(i) {
			// Reply "it is a private chat"
			handlerPack.NoticePrivateChat(s, i, b.FMTer)
			return
		}
		log.Printf("Executed guild ID: %v", i.GuildID)

		// If command use not use for game interaction
		if !cmd.IsUsedForGame() {
			// Just execute a Execute()
			log.Printf("Execute %v command.", cmdName)
			cmd.Execute(ctx, i.Interaction, nil)
			return
		}

		// I know the guildID
		executedGuildID := i.GuildID
		registerNewGame := func() {
			log.Printf("Must be register_game: Execute %v command.", cmdName)

			userRenameProvider := userPack.NewBotUserRenameProvider(s, executedGuildID)
			gameConfig := botGamePack.GetNewGameConfig(userRenameProvider, b.Databases.Storage)

			b.Games[executedGuildID] = gamePack.GetNewGame(executedGuildID, gameConfig...)
			content, isOk := handlerPack.ValidateCommandByGameState(
				s, executedCommandName, b.Games[executedGuildID], b.FMTer, b.Databases)
			if !isOk {
				handlerPack.Response(s, i.Interaction, content)
				return
			}
			log.Printf("Registered new game by %v Guild ID", executedGuildID)
			cmd.Execute(ctx, i.Interaction, b.Games[executedGuildID])
		}
		// And is there a game on this server
		_, containsGame := b.Games[executedGuildID]

		// If yes
		if containsGame {
			log.Printf("Execute %v command.", cmdName)
			currGame := b.Games[executedGuildID]
			if currGame.GetState() == gamePack.FinishState {
				registerNewGame()
				return
			}
			// Validate Is correct command by game state
			content, isOk := handlerPack.ValidateCommandByGameState(
				s, executedCommandName, currGame, b.FMTer, b.Databases)
			if !isOk {
				handlerPack.Response(s, i.Interaction, content)
				return
			}
			// If ok, I call the Execute method of the command
			cmd.Execute(ctx, i.Interaction, currGame)
			// If command is (finishGame), delete game from map
			if executedCommandName == handlerPack.FinishGameCommandName {
				b.Lock()
				delete(b.Games, i.GuildID)
				b.Unlock()
			}
			return
		}

		// Otherwise I know the game isn't registered.
		// I check to see if the command name is register_game. If not, it means that the
		// person uses the game command without registering it.
		if executedCommandName != handlerPack.RegisterGameCommandName {
			handlerPack.NoticeIsEmptyGame(s, i, b.FMTer)
			return
		}

		// Stand, that game is nil.
		// RegisterNewGame
		registerNewGame()
	}
}

func (b *Bot) registerCommands() {
	log.Println("Register commands")
	stateId := b.Session.State.User.ID
	for _, cmd := range b.Commands {
		newCmd := cmd
		registeredCmd, err := b.Session.ApplicationCommandCreate(stateId, "", newCmd.GetCmd())
		if err != nil {
			log.Print(err)
		}
		b.registeredCommands = append(b.registeredCommands, registeredCmd)
		log.Printf("Registered command, name %v", registeredCmd.Name)
	}
}

func (b *Bot) removeRegisteredCommands() { b.deleteAllGloballyRegisteredCommands() }

// deleteAllGloballyRegisteredCommands Delete all registered to bot functions. Globally Registered
func (b *Bot) deleteAllGloballyRegisteredCommands() {
	log.Println("Init Delete all globally registered commands.")
	userId := b.Session.State.User.ID
	globallyRegisteredCommands, err := b.Session.ApplicationCommands(userId, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, command := range globallyRegisteredCommands {
		log.Printf("Removed command ID: %v", command.ID)
		err = b.Session.ApplicationCommandDelete(userId, "", command.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("All global commands deleted.")
}

func (b *Bot) deleteHandlers() {
	err := b.Session.AddHandler(nil)
	if err != nil { // (Code by ChatGPT <3, lol)
		log.Println("Delete all handlers")
		return
	}
}

func (b *Bot) finishAllGames() {
	for _, game := range b.Games {
		ch := make(chan gamePack.Signal)
		go game.FinishAnyway()

		for {
			gSignal, ok := <-ch
			if !ok {
				break
			}
			log.Println(gSignal)
		}
	}
}
