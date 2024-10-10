package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/https-whoyan/MafiaBot/pkg"

	fmtErPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	namesPack "github.com/https-whoyan/MafiaBot/internal/handlers/names"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"
	coreRolesPack "github.com/https-whoyan/MafiaCore/roles"

	"github.com/bwmarrin/discordgo"
)

// _________________________________
// ********************************
// This contains all bot commands.
//********************************
// _________________________________

// _______________________
// Channels
// _______________________

var (
	loadLoggersOnce  = &sync.Once{}
	basicInfoLogger  = log.New(os.Stdout, "Info\t", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	basicErrorLogger = log.New(os.Stdout, "Error\t", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

func LoadLoggers(infoLogger, errLogger *log.Logger) {
	loadLoggersOnce.Do(func() {
		basicInfoLogger = infoLogger
		basicErrorLogger = errLogger
	})
}

type basicCmd struct {
	s             *discordgo.Session
	f             *fmtErPack.DiscordFMTer
	cmd           *discordgo.ApplicationCommand
	db            *pkg.Database
	isUsedForGame bool
	name          string
	infoLogger    *log.Logger
	errLogger     *log.Logger
}

func newBasicCmd(
	s *discordgo.Session, cmd *discordgo.ApplicationCommand, db *pkg.Database,
	name string, isUsedForGame bool,
) basicCmd {
	return basicCmd{
		s:             s,
		f:             fmtErPack.DiscordFMTInstance,
		cmd:           cmd,
		name:          name,
		isUsedForGame: isUsedForGame,
		db:            db,
		infoLogger: log.New(
			basicInfoLogger.Writer(),
			basicInfoLogger.Prefix()+"CommandName: "+name+"\t",
			basicInfoLogger.Flags(),
		),
		errLogger: log.New(
			basicErrorLogger.Writer(),
			basicErrorLogger.Prefix()+"CommandName: "+name+"\t",
			basicErrorLogger.Flags(),
		),
	}
}

func (c basicCmd) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c basicCmd) GetName() string                       { return c.name }
func (c basicCmd) IsUsedForGame() bool                   { return c.isUsedForGame }
func (c basicCmd) Execute(_ context.Context, _ *discordgo.Interaction, _ *coreGamePack.Game) {
	panic("implement me")
}

func (c basicCmd) log(i *discordgo.Interaction, g *coreGamePack.Game) {
	var (
		gameState       = "nil"
		gameNightVoting = "nil"
	)
	if g != nil {
		gameState = g.GetState().String()
		gameNightVoting = g.GetNightVoting().Name
	}
	message := fmt.Sprintf(
		"Execute %v command in %v GuildID. GameState: %v, GameNightVoting: %v",
		i.ApplicationCommandData().Name,
		i.GuildID,
		gameState,
		gameNightVoting,
	)
	c.infoLogger.Println(message)
}

func (c basicCmd) response(i *discordgo.Interaction, content string) {
	err := c.s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		c.errLogger.Print("Error send response: ", err)
	}
}

// addChannelRoleCommand command logic
type addChannelRoleCommand struct {
	basicCmd
}

func NewAddChannelRoleCommand(s *discordgo.Session, db *pkg.Database) Command {
	generateOption := func(roleName string) *discordgo.ApplicationCommandOption {
		return &discordgo.ApplicationCommandOption{
			Name:        roleName,
			Description: fmt.Sprintf("Add %s interationChat", roleName),
			Type:        discordgo.ApplicationCommandOptionString,
		}
	}

	generateOptions := func() []*discordgo.ApplicationCommandOption {
		allNamesOfRoles := coreRolesPack.GetInteractionRoleNamesWhoHasOwnChat()
		var options []*discordgo.ApplicationCommandOption
		for _, roleName := range allNamesOfRoles {
			options = append(options, generateOption(strings.ToLower(roleName)))
		}
		return options
	}

	return &addChannelRoleCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.AddChannelRoleCommandName,
				Description: "Define a chat room where the interaction between the bot and the role will take place.",
				Options:     generateOptions(),
			},
			db,
			namesPack.AddChannelRoleCommandName,
			false,
		),
	}
}

// addMainChannelCommand command logic
type addMainChannelCommand struct {
	basicCmd
}

func NewAddMainChannelCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &addMainChannelCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.AddMainChannelCommandName,
				Description: "Define a chat room where the interaction between the bot and all game participants.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "chat_id",
						Description: "Add main game chat.",
						Required:    true,
						Type:        discordgo.ApplicationCommandOptionString,
					},
				},
			},
			db,
			namesPack.AddMainChannelCommandName,
			false,
		),
	}
}

// _____________________
// Game
// _____________________

// registerGameCommand command logic
type registerGameCommand struct {
	basicCmd
}

func NewRegisterGameCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &registerGameCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.RegisterGameCommandName,
				Description: "Register new Game",
			},
			db,
			namesPack.RegisterGameCommandName,
			true,
		),
	}
}

type finishGameCommand struct {
	basicCmd
}

func NewFinishGameCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &finishGameCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.FinishGameCommandName,
				Description: "Ends the game early.",
			},
			db,
			namesPack.FinishGameCommandName,
			true,
		),
	}
}

// choiceGameConfigCommand command logic
type choiceGameConfigCommand struct {
	basicCmd
}

func NewChoiceGameConfigCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &choiceGameConfigCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.ChoiceGameConfigCommandName,
				Description: "This output a list of game configs for voting.",
			},
			db,
			namesPack.ChoiceGameConfigCommandName,
			true,
		),
	}
}

type startGameCommand struct {
	basicCmd
}

func NewStartGameCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &startGameCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.StartGameCommandName,
				Description: "Init game after game config choosing",
			},
			db,
			namesPack.StartGameCommandName,
			true,
		),
	}
}

// ______________
// Voting
// ______________

type gameVoteCommand struct {
	basicCmd
}

func NewGameVoteCommand(s *discordgo.Session, db *pkg.Database) Command {
	description := "The command used for voting. Put " + coreGamePack.EmptyVoteStr + " for empty vote."
	return &gameVoteCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.VoteGameCommandName,
				Description: description,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "player_id",
						Description: "Enter the player's game ID",
						Required:    true,
					},
				},
			},
			db,
			namesPack.VoteGameCommandName,
			true,
		),
	}
}

type gameTwoVoteCommand struct {
	basicCmd
}

func NewGameTwoVoteCommand(s *discordgo.Session, db *pkg.Database) Command {
	description := "The command used for voting, but only for roles that use 2 voices at once."
	return &gameTwoVoteCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.TwoVoteGameCommandName,
				Description: description,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "player_id_1",
						Description: "Enter the player's game ID",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "player_id_2",
						Description: "Enter the player's game ID",
						Required:    true,
					},
				},
			},
			db,
			namesPack.TwoVoteGameCommandName,
			true,
		),
	}
}

type dayVoteCommand struct {
	basicCmd
}

func NewDayVoteCommand(s *discordgo.Session, db *pkg.Database) Command {
	description := "The command used for day voting, use " + coreGamePack.EmptyVoteStr + " for skip."
	return &dayVoteCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.DayVoteGameCommandName,
				Description: description,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "kicked_player_id",
						Description: "Enter the player's game ID",
						Required:    true,
					},
				},
			},
			db,
			namesPack.DayVoteGameCommandName,
			true,
		),
	}
}

// ___________
// Other
// ___________

// yanLohCommand command
type yanLohCommand struct {
	basicCmd
}

func NewYanLohCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &yanLohCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.YanLohCommandName,
				Description: "Call Yan with this command!",
			},
			db,
			namesPack.YanLohCommandName,
			false,
		),
	}
}

// aboutRolesCommand command logic
type aboutRolesCommand struct {
	basicCmd
}

func NewAboutRolesCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &aboutRolesCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        namesPack.AboutRolesCommandName,
				Description: "Send description about roles",
			},
			db,
			namesPack.AboutRolesCommandName,
			false,
		),
	}
}
