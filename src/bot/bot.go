package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/pkg/channel"
	"github.com/https-whoyan/MafiaBot/pkg/game"
	"log"
	"os"
)

type Bot struct {
	token              string
	Session            *discordgo.Session
	Commands           map[string]Command
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
	commandName := c.GetName()
	b.Commands[commandName] = c
}

func (b *Bot) InitBotCommands() {
	b.initCommand(NewYanLohCommand())
	b.initCommand(channel.NewAddChannelRole())
	b.initCommand(game.NewRegisterGameCommand())
}

func (b *Bot) RegisterHandlers() {
	log.Print("Register handlers")
	for _, cmd := range b.Commands {
		cmdName := cmd.GetName()
		newCmd := cmd.GetExecuteFunc()
		log.Printf("Register handler, command name: %v, func hash: %v", cmdName, newCmd)
		b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if ok := cmdName == i.ApplicationCommandData().Name; ok {
				log.Printf("Execute %v command.", cmdName)
				newCmd(s, i)
			}
		})
	}
}

func (b *Bot) RegisterCommands() {
	log.Println("Register commands")
	stateId := b.Session.State.User.ID
	for _, cmd := range b.Commands {
		registeredCmd, err := b.Session.ApplicationCommandCreate(stateId, "", cmd.GetCmd())
		if err != nil {
			log.Print(err)
		}
		log.Printf("Register command, name %v", registeredCmd.Name)
		b.registeredCommands = append(b.registeredCommands, registeredCmd)
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
