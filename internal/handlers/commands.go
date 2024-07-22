package handlers

import (
	"fmt"
	"strings"

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
	AddChannelRoleCommandName   = "add_channel_role"
	AddMainChannelCommandName   = "add_main_channel"
	RegisterGameCommandName     = "register_game"
	ChoiceGameConfigCommandName = "choose_game_config"
	YanLohCommandName           = "yan_loh"
	AboutRolesCommandName       = "about_roles"
	StartGameCommandName        = "start_game"
	VoteGameCommandName         = "vote"
	TwoVoteGameCommandName      = "two_vote"
	DayVoteGameCommandName      = "day_vote"
	FinishGameCommandName       = "finish_game"
)

// _______________________
// Channels
// _______________________

// AddChannelRoleCommand command logic
type AddChannelRoleCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAddChannelRoleCommand() *AddChannelRoleCommand {
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

	return &AddChannelRoleCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        AddChannelRoleCommandName,
			Description: "Define a chat room where the interaction between the bot and the role will take place.",
			Options:     generateOptions(),
		},
		isUsedForGame: false,
		name:          AddChannelRoleCommandName,
	}
}

func (c AddChannelRoleCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c AddChannelRoleCommand) GetName() string                       { return c.name }
func (c AddChannelRoleCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

// AddMainChannelCommand command logic
type AddMainChannelCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAddMainChannelCommand() *AddMainChannelCommand {
	return &AddMainChannelCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        AddMainChannelCommandName,
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
		isUsedForGame: false,
		name:          AddMainChannelCommandName,
	}
}

func (c AddMainChannelCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c AddMainChannelCommand) GetName() string                       { return c.name }
func (c AddMainChannelCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

// _____________________
// Game
// _____________________

// RegisterGameCommand command logic
type RegisterGameCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewRegisterGameCommand() *RegisterGameCommand {
	return &RegisterGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        RegisterGameCommandName,
			Description: "Register new Game",
		},
		isUsedForGame: true,
		name:          RegisterGameCommandName,
	}
}

func (c RegisterGameCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c RegisterGameCommand) GetName() string                       { return c.name }
func (c RegisterGameCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

type FinishGameCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewFinishGameCommand() *FinishGameCommand {
	return &FinishGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        FinishGameCommandName,
			Description: "Ends the game early.",
		},
		isUsedForGame: true,
		name:          FinishGameCommandName,
	}
}

func (c FinishGameCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c FinishGameCommand) GetName() string                       { return c.name }
func (c FinishGameCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

// ChoiceGameConfigCommand command logic
type ChoiceGameConfigCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewChoiceGameConfigCommand() *ChoiceGameConfigCommand {
	return &ChoiceGameConfigCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        ChoiceGameConfigCommandName,
			Description: "This output a list of game configs for voting.",
		},
		isUsedForGame: true,
		name:          ChoiceGameConfigCommandName,
	}
}

func (c ChoiceGameConfigCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c ChoiceGameConfigCommand) GetName() string                       { return c.name }
func (c ChoiceGameConfigCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

type StartGameCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewStartGameCommand() *StartGameCommand {
	return &StartGameCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        StartGameCommandName,
			Description: "Init game after game config choosing",
		},
		isUsedForGame: true,
		name:          StartGameCommandName,
	}
}

func (c StartGameCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c StartGameCommand) GetName() string                       { return c.name }
func (c StartGameCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

// ______________
// Voting
// ______________

type GameVoteCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewGameVoteCommand() *GameVoteCommand {
	description := "The command used for voting. Put " + coreGamePack.EmptyVoteStr + " for empty vote."
	return &GameVoteCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        VoteGameCommandName,
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
		isUsedForGame: true,
		name:          VoteGameCommandName,
	}
}

func (c GameVoteCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c GameVoteCommand) GetName() string                       { return c.name }
func (c GameVoteCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

type GameTwoVoteCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewGameTwoVoteCommand() *GameTwoVoteCommand {
	description := "The command used for voting, but only for roles that use 2 voices at once."
	return &GameTwoVoteCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        TwoVoteGameCommandName,
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
		isUsedForGame: true,
		name:          TwoVoteGameCommandName,
	}
}

func (c GameTwoVoteCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c GameTwoVoteCommand) GetName() string                       { return c.name }
func (c GameTwoVoteCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

type DayVoteCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewDayVoteCommand() *DayVoteCommand {
	description := "The command used for day voting, use " + coreGamePack.EmptyVoteStr + " for skip."
	return &DayVoteCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        DayVoteGameCommandName,
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
		isUsedForGame: true,
		name:          DayVoteGameCommandName,
	}
}

func (c DayVoteCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c DayVoteCommand) GetName() string                       { return c.name }
func (c DayVoteCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

// ___________
// Other
// ___________

// YanLohCommand command
type YanLohCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewYanLohCommand() *YanLohCommand {
	return &YanLohCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        YanLohCommandName,
			Description: "Call Yan with this command!",
		},
		isUsedForGame: false,
		name:          YanLohCommandName,
	}
}

func (c YanLohCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c YanLohCommand) GetName() string                       { return c.name }
func (c YanLohCommand) IsUsedForGame() bool                   { return c.isUsedForGame }

// AboutRolesCommand command logic
type AboutRolesCommand struct {
	cmd           *discordgo.ApplicationCommand
	isUsedForGame bool
	name          string
}

func NewAboutRolesCommand() *AboutRolesCommand {
	return &AboutRolesCommand{
		cmd: &discordgo.ApplicationCommand{
			Name:        AboutRolesCommandName,
			Description: "Send description about roles",
		},
		isUsedForGame: false,
		name:          AboutRolesCommandName,
	}
}

func (c AboutRolesCommand) GetCmd() *discordgo.ApplicationCommand { return c.cmd }
func (c AboutRolesCommand) GetName() string                       { return c.name }
func (c AboutRolesCommand) IsUsedForGame() bool                   { return c.isUsedForGame }
