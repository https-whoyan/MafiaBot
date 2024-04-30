package bot

import (
	"fmt"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"

	"github.com/https-whoyan/MafiaBot/internal/bot/channel"
	botGame "github.com/https-whoyan/MafiaBot/internal/bot/game/registration"
)

type Bot struct {
	sync.RWMutex
	token              string
	Session            *discordgo.Session
	Commands           map[string]Command
	Game               *game.Game
	registeredCommands []*discordgo.ApplicationCommand
}

func InitBot() *Bot {
	token := os.Getenv("BOT_TOKEN")
	botStr := "Bot " + token
	s, err := discordgo.New(botStr)
	if err != nil {
		log.Fatal(err)
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	return &Bot{
		token:    token,
		Session:  s,
		Commands: make(map[string]Command),
		Game:     game.NewUndefinedGame(),
	}
}
func (b *Bot) Close() {
	err := b.Session.Close()
	if err != nil {
		log.Printf("err in closing bot, err: %v", err)
	}
}

func (b *Bot) HasACommands() bool {
	userId := b.Session.State.User.ID
	globallyRegisteredCommands, _ := b.Session.ApplicationCommands(userId, "")
	return len(globallyRegisteredCommands) == 0
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

func (b *Bot) initCommand(c Command) {
	b.Lock()
	defer b.Unlock()
	commandName := c.GetName()
	b.Commands[commandName] = c
}

func (b *Bot) InitBotCommands() {
	b.initCommand(NewYanLohCommand())
	b.initCommand(channel.NewAddChannelRole())
	b.initCommand(botGame.NewRegisterGameCommand())
}

func (b *Bot) RegisterHandlers() {
	log.Print("Register handlers")
	for _, cmd := range b.Commands {
		newCmd := cmd
		cmdName := newCmd.GetName()
		log.Printf("Register handler, command name: %v", cmdName)
		b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("я тут")
			executedCommandName := i.ApplicationCommandData().Name
			//executedCommandIsInPrivateChat := false
			if cmdName == executedCommandName {
				log.Printf("Execute %v command.", cmdName)
				newCmd.Execute(s, i)
				b.Game.Lock()
				log.Printf("Exectute %v game interation podcommand", cmdName)
				b.Game = newCmd.GameInteraction(b.Game)
				b.Game.Unlock()
			}
		})
	}

}

func (b *Bot) RegisterCommands() {
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

func (b *Bot) RemoveRegisteredCommands() {
	log.Println("Remove commands")
	stateId := b.Session.State.User.ID
	for _, registeredCmd := range b.registeredCommands {
		err := b.Session.ApplicationCommandDelete(stateId, "", registeredCmd.ID)
		if err != nil {
			log.Printf("cannot delete command, err: %v", err)
		}
	}
}
