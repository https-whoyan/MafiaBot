package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/https-whoyan/MafiaBot/pkg"
	"github.com/https-whoyan/MafiaBot/pkg/repository/mongo"
	"github.com/https-whoyan/MafiaBot/pkg/repository/redis"

	fmtErPack "github.com/https-whoyan/MafiaBot/internal/fmt"
	coreGamePack "github.com/https-whoyan/MafiaCore/game"
	coreRolesPack "github.com/https-whoyan/MafiaCore/roles"

	"github.com/bwmarrin/discordgo"
)

// _________________________________
// ********************************
// This contains all bot commands.
//********************************
// _________________________________

const (
	addChannelRoleCommandName   = "add_channel_role"
	addMainChannelCommandName   = "add_main_channel"
	RegisterGameCommandName     = "register_game"
	choiceGameConfigCommandName = "choose_game_config"
	yanLohCommandName           = "yan_loh"
	aboutRolesCommandName       = "about_roles"
	startGameCommandName        = "start_game"
	voteGameCommandName         = "vote"
	twoVoteGameCommandName      = "two_vote"
	dayVoteGameCommandName      = "day_vote"
	FinishGameCommandName       = "finish_game"
)

// _______________________
// Channels
// _______________________

type basicCmd struct {
	s             *discordgo.Session
	f             *fmtErPack.DiscordFMTer
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func newBasicCmd(s *discordgo.Session, cmd *discordgo.ApplicationCommand, name string, isUsedForGame bool) basicCmd {
	return basicCmd{
		s:             s,
		f:             fmtErPack.DiscordFMTInstance,
		cmd:           cmd,
		name:          name,
		isUsedForGame: isUsedForGame,
	}
}

func (c basicCmd) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c basicCmd) GetName() string                       { return c.name }
func (c basicCmd) IsUsedForGame() bool                   { return c.isUsedForGame }
func (c basicCmd) Execute(_ context.Context, _ *discordgo.Interaction, _ *coreGamePack.Game) {
	panic("implement me")
}

func (c basicCmd) response(i *discordgo.Interaction, content string) {
	err := c.s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Print(err)
	}
}

// addChannelRoleCommand command logic
type addChannelRoleCommand struct {
	basicCmd
	storage mongo.Storage
}

func NewAddChannelRoleCommand(s *discordgo.Session, storage mongo.Storage) Command {
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
				Name:        addChannelRoleCommandName,
				Description: "Define a chat room where the interaction between the bot and the role will take place.",
				Options:     generateOptions(),
			},
			addChannelRoleCommandName,
			false,
		),
		storage: storage,
	}
}

// addMainChannelCommand command logic
type addMainChannelCommand struct {
	basicCmd
	storage mongo.Storage
}

func NewAddMainChannelCommand(s *discordgo.Session, storage mongo.Storage) Command {
	return &addMainChannelCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        addMainChannelCommandName,
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
			addMainChannelCommandName,
			false,
		),
		storage: storage,
	}
}

// _____________________
// Game
// _____________________

// registerGameCommand command logic
type registerGameCommand struct {
	basicCmd
	db *pkg.Database
}

func NewRegisterGameCommand(s *discordgo.Session, db *pkg.Database) Command {
	return &registerGameCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        RegisterGameCommandName,
				Description: "Register new Game",
			},
			RegisterGameCommandName,
			false,
		),
		db: db,
	}
}

type finishGameCommand struct {
	basicCmd
}

func NewFinishGameCommand(s *discordgo.Session) Command {
	return &finishGameCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        FinishGameCommandName,
				Description: "Ends the game early.",
			},
			FinishGameCommandName,
			true,
		),
	}
}

// choiceGameConfigCommand command logic
type choiceGameConfigCommand struct {
	basicCmd
	hasher redis.Hasher
}

func NewChoiceGameConfigCommand(s *discordgo.Session, hasher redis.Hasher) Command {
	return &choiceGameConfigCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        choiceGameConfigCommandName,
				Description: "This output a list of game configs for voting.",
			},
			choiceGameConfigCommandName,
			true,
		),
		hasher: hasher,
	}
}

type startGameCommand struct {
	basicCmd
	hasher redis.Hasher
}

func NewStartGameCommand(s *discordgo.Session, hasher redis.Hasher) Command {
	return &startGameCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        startGameCommandName,
				Description: "Init game after game config choosing",
			},
			startGameCommandName,
			true,
		),
		hasher: hasher,
	}
}

// ______________
// Voting
// ______________

type gameVoteCommand struct {
	basicCmd
}

func NewGameVoteCommand(s *discordgo.Session) Command {
	description := "The command used for voting. Put " + coreGamePack.EmptyVoteStr + " for empty vote."
	return &gameVoteCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        voteGameCommandName,
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
			voteGameCommandName,
			true,
		),
	}
}

type gameTwoVoteCommand struct {
	basicCmd
}

func NewGameTwoVoteCommand(s *discordgo.Session) Command {
	description := "The command used for voting, but only for roles that use 2 voices at once."
	return &gameTwoVoteCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        twoVoteGameCommandName,
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
			twoVoteGameCommandName,
			true,
		),
	}
}

type dayVoteCommand struct {
	basicCmd
}

func NewDayVoteCommand(s *discordgo.Session) Command {
	description := "The command used for day voting, use " + coreGamePack.EmptyVoteStr + " for skip."
	return &dayVoteCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        dayVoteGameCommandName,
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
			dayVoteGameCommandName,
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

func NewYanLohCommand(s *discordgo.Session) Command {
	return &yanLohCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        yanLohCommandName,
				Description: "Call Yan with this command!",
			},
			yanLohCommandName,
			false,
		),
	}
}

// aboutRolesCommand command logic
type aboutRolesCommand struct {
	basicCmd
}

func NewAboutRolesCommand(s *discordgo.Session) Command {
	return &aboutRolesCommand{
		basicCmd: newBasicCmd(
			s,
			&discordgo.ApplicationCommand{
				Name:        aboutRolesCommandName,
				Description: "Send description about roles",
			},
			aboutRolesCommandName,
			false,
		),
	}
}
