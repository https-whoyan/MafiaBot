package bot

import (
	"errors"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	Commands map[string]Command
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
			Commands: make(map[string]Command),
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

func (b *Bot) initCommand(c Command) {
	b.Lock()
	defer b.Unlock()
	commandName := c.GetName()
	b.Commands[commandName] = c
}

func (b *Bot) initBotCommands() {
	b.initCommand(NewYanLohCommand())
	b.initCommand(NewAddChannelRole())
	b.initCommand(NewRegisterGameCommand())
}

func (b *Bot) registerHandlers() {
	log.Print("Register handlers")
	for _, cmd := range b.Commands {
		newCmd := cmd
		cmdName := newCmd.GetName()
		log.Printf("Register handler, command name: %v", cmdName)
		b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			executedCommandName := i.ApplicationCommandData().Name

			if executedCommandName != cmdName {
				return
			}

			if isPrivateMessage(i) {
				noticePrivateChat(s, i)
				return
			}

			executedGuildID := i.GuildID
			currGame, containsGame := b.Games[executedGuildID]

			log.Printf("Execute %v command.", cmdName)
			newCmd.Execute(s, i.Interaction)

			if containsGame {
				log.Printf("Execute %v game interation podcommand", cmdName)
				currGame.Lock()
				newCmd.GameInteraction(currGame)
				currGame.Unlock()
				return
			}
			if executedCommandName != "register_game" {
				noticeIsEmptyGame(s, i)
				return
			}
			currGame = &game.Game{}
			newCmd.GameInteraction(currGame)
			return
		})
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
