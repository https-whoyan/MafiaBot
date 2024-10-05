package bot

import (
	"context"
	"github.com/https-whoyan/MafiaBot/internal/game"
	"github.com/https-whoyan/MafiaBot/internal/handlers/names"
	userPack "github.com/https-whoyan/MafiaBot/internal/user"
	"log"
	"os"
	"sync"

	botFMTPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	handlerPack "github.com/https-whoyan/MafiaBot/internal/handlers"
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
	// Implement of FmtInterface.
	FMTer      *botFMTPack.DiscordFMTer
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
	// Databases
	Databases *pkgPack.Database
}

func InitBot(ctx context.Context, cnf *Config, databases *pkgPack.Database,
	errLogger *log.Logger, infoLogger *log.Logger) (*Bot, error) {
	token := cnf.token
	botStr := "Bot " + token
	s, err := discordgo.New(botStr)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		token:      token,
		Session:    s,
		Commands:   make(map[string]handlerPack.Command),
		Games:      make(map[string]*gamePack.Game),
		FMTer:      botFMTPack.DiscordFMTInstance,
		InfoLogger: infoLogger,
		ErrLogger:  errLogger,
		Databases:  databases,
	}
	handlerPack.LoadLoggers(infoLogger, errLogger)

	bot.initBotCommands()
	bot.registerHandlers(ctx)
	return bot, nil
}

func (b *Bot) loginAs() {
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.InfoLogger.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
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
	b.initCommand(handlerPack.NewAddMainChannelCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewAddChannelRoleCommand(b.Session, b.Databases))

	// Game
	b.initCommand(handlerPack.NewRegisterGameCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewChoiceGameConfigCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewStartGameCommand(b.Session, b.Databases))

	// Vote
	b.initCommand(handlerPack.NewGameVoteCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewGameTwoVoteCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewDayVoteCommand(b.Session, b.Databases))

	// finish game
	b.initCommand(handlerPack.NewFinishGameCommand(b.Session, b.Databases))

	// Other
	b.initCommand(handlerPack.NewYanLohCommand(b.Session, b.Databases))
	b.initCommand(handlerPack.NewAboutRolesCommand(b.Session, b.Databases))

}

func (b *Bot) registerHandlers(ctx context.Context) {
	b.InfoLogger.Print("Register handlers")
	for _, cmd := range b.Commands {
		// To avoid closing the loop
		newCmd := cmd
		cmdName := newCmd.GetName()
		b.InfoLogger.Printf("Register handler, command name: %v", cmdName)
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
				b.ErrLogger.Printf("Recovered from panic: %v", r)
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
			handlerPack.NoticePrivateChat(s, i, b.FMTer, b.ErrLogger)
			return
		}
		// If game not used for change game
		if !cmd.IsUsedForGame() {
			// Execute it
			cmd.Execute(ctx, i.Interaction, nil)
			return
		}
		// Check, contains a game
		currGame, isContains := b.Games[i.GuildID]
		if !isContains {
			// Maybe it a registration game command
			if cmdName != names.RegisterGameCommandName {
				handlerPack.NoticeIsEmptyGame(s, i, b.FMTer, b.ErrLogger)
				return
			}
			gameOptions := game.GetNewGameConfig(
				userPack.NewBotUserRenameProvider(s, i.GuildID),
				b.Databases.Storage,
				b.ErrLogger,
				b.InfoLogger,
			)
			newGame := gamePack.GetNewGame(i.GuildID, gameOptions...)
			b.Lock()
			b.Games[i.GuildID] = newGame
			b.Unlock()
			cmd.Execute(ctx, i.Interaction, newGame)
		}
		// Validate the stage
		content, ok := handlerPack.ValidateCommandByGameState(
			s, b.Databases, currGame, cmdName, b.FMTer,
		)
		if !ok {
			handlerPack.Response(s, i.Interaction, content, b.ErrLogger)
			return
		}
		// Ok, execute
		cmd.Execute(ctx, i.Interaction, currGame)
	}
}

func (b *Bot) registerCommands() {
	b.InfoLogger.Println("Register commands")
	stateId := b.Session.State.User.ID
	for _, cmd := range b.Commands {
		newCmd := cmd
		registeredCmd, err := b.Session.ApplicationCommandCreate(stateId, "", newCmd.GetCmd())
		if err != nil {
			b.ErrLogger.Print(err)
		}
		b.registeredCommands = append(b.registeredCommands, registeredCmd)
		b.InfoLogger.Printf("Registered command, name %v", registeredCmd.Name)
	}
}

func (b *Bot) removeRegisteredCommands() { b.deleteAllGloballyRegisteredCommands() }

// deleteAllGloballyRegisteredCommands Delete all registered to bot functions. Globally Registered
func (b *Bot) deleteAllGloballyRegisteredCommands() {
	b.InfoLogger.Println("Init Delete all globally registered commands.")
	userId := b.Session.State.User.ID
	globallyRegisteredCommands, err := b.Session.ApplicationCommands(userId, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, command := range globallyRegisteredCommands {
		b.InfoLogger.Printf("Removed command ID: %v", command.ID)
		err = b.Session.ApplicationCommandDelete(userId, "", command.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	b.InfoLogger.Println("All global commands deleted.")
}

func (b *Bot) deleteHandlers() {
	err := b.Session.AddHandler(nil)
	if err != nil { // (Code by ChatGPT <3, lol)
		b.InfoLogger.Println("Delete all handlers")
		return
	}
}

func (b *Bot) finishAllGames() {
	for _, runningGame := range b.Games {
		b.InfoLogger.Printf("Finish game %v", runningGame.GuildID)
		runningGame.FinishAnyway()
	}
}
