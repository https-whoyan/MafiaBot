package bot

import (
	"errors"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
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

func InitBot(token string) (*Bot, error) {
	botStr := "Bot " + token
	s, err := discordgo.New(botStr)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		token:    token,
		Session:  s,
		Commands: make(map[string]Command),
		Games:    make(map[string]*game.Game),
	}
	bot.initBotCommands()
	bot.registerHandlers()
	err = bot.Open()
	if err != nil {
		return nil, err
	}
	bot.loginAs()
	bot.registerCommands()
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
	return nil
}

func (b *Bot) Close() {
	b.DeleteHandlers()
	err := b.Session.Close()
	if err != nil {
		log.Printf("Error clothing bot, err: %v", err)
	}
	b.removeRegisteredCommands()

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
			executedCommandIsInPrivateChat := i.GuildID == ""
			if executedCommandIsInPrivateChat && executedCommandName == cmdName {
				content := "All commands are used on the server. If you have difficulties in using the bot, " +
					"please refer to the repository documentation: https://github.com/https-whoyan/MafiaBot."
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
				if err != nil {
					log.Println(errors.Join(
						errors.New("there was an error when sending a private message, err: "), err),
					)
				}
				return
			}

			executedGuildID := i.GuildID
			currGame, containsGame := b.Games[executedGuildID]

			if cmdName == executedCommandName {
				log.Printf("Execute %v command.", cmdName)
				newCmd.Execute(s, i.Interaction)

				if containsGame {
					log.Printf("Execute %v game interation podcommand", cmdName)
					currGame.Lock()
					newCmd.GameInteraction(currGame)
					currGame.Unlock()
				} else {
					if executedCommandName == "register_game" {
						currGame = &game.Game{}
						newCmd.GameInteraction(currGame)
					} else {
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You can't interact with the game because you haven't started" +
									" it. Write the /register_game command to start the game.",
							},
						})
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
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
	log.Printf("Break program.")
	os.Exit(0)
}

func (b *Bot) DeleteHandlers() {
	err := b.Session.AddHandler(nil)
	if err != nil { // :))))
		log.Println("Delete all handlers")
		return
	}
}
