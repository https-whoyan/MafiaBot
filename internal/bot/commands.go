package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/game"
	time2 "github.com/https-whoyan/MafiaBot/internal/time"
	"github.com/https-whoyan/MafiaBot/pkg/db/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/db/redis"
	"log"
	"strconv"
	"time"
)

// Registration variables
var (
	RegistrationPlayerSticker    = ":grin:"
	RegistrationSpectatorSticker = ":smiling_imp:"
)

// YanLohCommand command
type YanLohCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewYanLohCommand() *YanLohCommand {
	name := "yan_loh"
	return &YanLohCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Call Yan with this command!",
		},
		name: name,
	}
}

func (c *YanLohCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *YanLohCommand) GetName() string {
	return c.name
}

func (c *YanLohCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {
	messageContent := "Возможно, что ян и лох. И древлян. Но что бы его же ботом его обзывать..."
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: messageContent,
		},
	})
	if err != nil {
		log.Print(err)
	}
	// async kick requester
	guildID := i.GuildID
	go func(sessId, kickedUserID string) {
		var kickPing time.Duration = 3
		time.Sleep(time.Second * kickPing)

		err = s.GuildMemberMove(sessId, kickedUserID, nil)
		if err != nil {
			log.Printf("failed kick user, err: %v", err)
		}
	}(guildID, i.Member.User.ID)
}

func (c *YanLohCommand) GameInteraction(g *game.Game) {
	//...
}

// RegisterGameCommand command
type RegisterGameCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewRegisterGameCommand() *RegisterGameCommand {
	name := "register_game"
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Register new Game",
		},
		name: name,
	}
}

func (c *RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *RegisterGameCommand) GetName() string {
	return c.name
}

func (c *RegisterGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ok. Message below.",
		},
	})
	if err != nil {
		log.Print(err)
	}
	deadlineStr := strconv.Itoa(time2.RegistrationDeadlineMinutes)
	responseMessageText := "Registration has begun. \n" +
		Bold("Post"+RegistrationPlayerSticker+" reactions below.") + Italic("If you want to be a spectator, "+
		"put the reaction "+RegistrationSpectatorSticker+".") + "\n\n" + Bold(
		Emphasized("Deadline: "+deadlineStr+"minutes</u>"))

	channelID := i.ChannelID
	message, err := s.ChannelMessageSend(channelID, responseMessageText)
	if err != nil {
		log.Print(err)
	}
	messageID := message.ID
	currDB, isContains := redis.GetCurrRedisDB()
	if !isContains {
		log.Print("DB isn't initialed, couldn't get currDB")
		return
	}
	err = currDB.SetInitialGameMessageID(i.GuildID, messageID)
	if err != nil {
		log.Print(err)
	}
}

func (c *RegisterGameCommand) GameInteraction(g *game.Game) {
	g = game.NewUndefinedGame()
}

// AddChannelRoleCommand command logic
type AddChannelRoleCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewAddChannelRole() *AddChannelRoleCommand {
	name := "add_channel_role"
	return &AddChannelRoleCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Define a chat room where the interaction between the bot and the role will take place.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "mafia",
					Description: "Add Mafia interaction chat",
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "doctor",
					Description: "Add Doctor interaction chat",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		name: name,
	}
}

func (c *AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *AddChannelRoleCommand) GetName() string {
	return c.name
}

func (c *AddChannelRoleCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {
	if len(i.ApplicationCommandData().Options) == 0 {
		err := s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Must be option!",
			},
		})
		if err != nil {
			log.Print(err)
		}
		return
	}

	roleName := i.ApplicationCommandData().Options[0].Name
	requestedChatID := i.ApplicationCommandData().Options[0].Value.(string)

	noticeChatContent := fmt.Sprintf("Now chat is used for %v role.", roleName)
	if _, ok := noticeChat(s, requestedChatID, noticeChatContent); ok != nil {
		err := s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid Chat ID!",
			},
		})
		if err != nil {
			log.Print(err)
		}
		return
	}

	messageContent := fmt.Sprintf("Done, now is %v chat is %v chat.", requestedChatID, roleName)
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: messageContent,
		},
	})
	if err != nil {
		log.Print(err)
	}

	currDB, containsDB := mongo.GetCurrMongoDB()
	if !containsDB {
		log.Println("empty database")
		return
	}
	currDB.Lock()
	currDB.Unlock()
	err = currDB.SetRoleChannel(i.ChannelID, requestedChatID, roleName)
	if err != nil {
		log.Println(err)
	}
	return
}

func (c *AddChannelRoleCommand) GameInteraction(g *game.Game) {
	//..
}

// StartGameCommand command logic
type StartGameCommand struct {
	cmd  *discordgo.ApplicationCommand
	name string
}

func NewStartGameCommand() *RegisterGameCommand {
	name := "start_game"
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        name,
			Description: "Start a new game. This output a list of game configs for voting",
		},
		name: name,
	}
}

func (c *StartGameCommand) GetCmd() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c *StartGameCommand) GetName() string {
	return c.name
}

func (c *StartGameCommand) Execute(s *discordgo.Session, i *discordgo.Interaction) {

}

func (c *StartGameCommand) GameInteraction(g *game.Game) {
	//..
}
