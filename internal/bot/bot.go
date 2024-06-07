package bot

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	h "github.com/https-whoyan/MafiaBot/internal/bot/handlers"
	"github.com/https-whoyan/MafiaBot/internal/core/game"

	"github.com/bwmarrin/discordgo"
)

// ____
// cfg
// ____

type BotConfig struct {
	token string
}

func LoadBotConfig() *BotConfig {
	token := os.Getenv("BOT_TOKEN")
	return &BotConfig{
		token: token,
	}
}

// ____
// Bot
// ____

var (
	botOnce     sync.Once
	botInstance *Bot
)

type Bot struct {
	sync.RWMutex
	token    string
	Session  *discordgo.Session
	Commands map[string]h.Command
	Games    map[string]*game.Game
	// Games this a map,
	// the key in which is the State of the server where the bot is running,
	// and the value is the game.
	registeredCommands []*discordgo.ApplicationCommand
}

func InitBot(cnf *BotConfig) {
	botOnce.Do(func() {
		token := cnf.token
		botStr := "Bot " + token
		s, err := discordgo.New(botStr)
		if err != nil {
			log.Fatal(err)
		}
		bot := &Bot{
			token:    token,
			Session:  s,
			Commands: make(map[string]h.Command),
			Games:    make(map[string]*game.Game),
		}
		bot.initBotCommands()
		bot.registerHandlers()
		botInstance = bot
	})
}

func Run() {
	if botInstance == nil {
		log.Fatal("Bot isn't instance!")
	}
	err := botInstance.Open()
	if err != nil {
		log.Fatal(err)
	}
	botInstance.loginAs()
	botInstance.registerCommands()

	// If you need delete all registered commands, use here: bot.DeleteAllGloballyRegisteredCommands()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func DisconnectBot() error {
	if botInstance == nil {
		return errors.New("bot isn't initialized")
	}
	return botInstance.Close()

}

func (b *Bot) loginAs() {
	b.Lock()
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	b.Unlock()
}

func (b *Bot) Open() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}
	b.loginAs()
	return nil
}

func (b *Bot) Close() error {
	b.Lock()
	defer b.Unlock()
	b.DeleteHandlers()
	err := b.Session.Close()
	if err != nil {
		return err
	}
	b.removeRegisteredCommands()
	return nil
}

func (b *Bot) initCommand(c h.Command) {
	b.Lock()
	defer b.Unlock()
	commandName := c.GetName()
	b.Commands[commandName] = c
}

func (b *Bot) initBotCommands() {
	b.initCommand(h.NewYanLohCommand())
	b.initCommand(h.NewAddChannelRole())
	b.initCommand(h.NewRegisterGameCommand())
	b.initCommand(h.NewChoiceGameConfig())
	b.initCommand(h.NewAboutRolesCommand())
}

func (b *Bot) registerHandlers() {
	log.Print("Register handlers")
	for _, cmd := range b.Commands {
		b.Lock()
		// To avoid closing the loop
		newCmd := cmd
		cmdName := newCmd.GetName()
		log.Printf("Register handler, command name: %v", cmdName)
		// Lock
		newHandler := b.getSIHandler(newCmd, cmdName)
		b.Session.AddHandler(newHandler)
		b.Unlock()
	}
}

func (b *Bot) getSIHandler(cmd h.Command, cmdName string) func(
	s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// I recognize the name of the team
		executedCommandName := i.ApplicationCommandData().Name

		// If it not equals as iterable cmd.Name
		if executedCommandName != cmdName {
			return
		}

		// If it executed in private chat
		if h.IsPrivateMessage(i) {
			// Reply "it is a private chat"
			h.NoticePrivateChat(s, i)
			return
		}

		// If command use not use for game interaction
		if !cmd.IsUsedForGame() {
			// Just execute a Execute()
			log.Printf("Execute %v command.", cmdName)
			cmd.Execute(s, i.Interaction, nil)
			return
		}

		// I know the guildID
		executedGuildID := i.GuildID
		log.Printf("Executed guild ID: %v", executedGuildID)
		// And is there a game on this server
		_, containsGame := b.Games[executedGuildID]

		// If yes
		if containsGame {
			log.Printf("Execute %v command.", cmdName)
			currGame := b.Games[executedGuildID]
			currGame.Lock()
			defer currGame.Unlock()
			// I call the Execute method of the command
			cmd.Execute(s, i.Interaction, currGame)
			return
		}

		// Otherwise I know the game isn't registered.
		//I check to see if the command name is register_game. If not, it means that the
		// person uses the game command without registering it.
		if executedCommandName != "register_game" {
			h.NoticeIsEmptyGame(s, i)
			return
		}

		// I use the register_game command
		log.Printf("Must be register_game: Execute %v command.", cmdName)
		b.Games[executedGuildID] = &game.Game{}
		currGame := b.Games[executedGuildID]
		currGame.Lock()
		defer currGame.Unlock()
		cmd.Execute(s, i.Interaction, currGame)

		return
	}
}

func (b *Bot) registerCommands() {
	log.Println("Register commands")
	b.Lock()
	defer b.Unlock()
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

func (b *Bot) removeRegisteredCommands() {
	b.DeleteAllGloballyRegisteredCommands()
}

// DeleteAllGloballyRegisteredCommands Delete all registered to bot functions. Globally Registered
func (b *Bot) DeleteAllGloballyRegisteredCommands() {
	log.Println("Start Delete all globally registered commands.")
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

func (b *Bot) DeleteHandlers() {
	err := b.Session.AddHandler(nil)
	if err != nil { // :))))
		log.Println("Delete all handlers")
		return
	}
}
