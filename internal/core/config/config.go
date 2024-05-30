package config

import (
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	"github.com/https-whoyan/MafiaBot/internal/core/roles/mafia"
	"github.com/https-whoyan/MafiaBot/internal/core/roles/maniac"
	"github.com/https-whoyan/MafiaBot/internal/core/roles/peaceful"
)

type RoleConfig struct {
	Role  *roles.Role `json:"role"`
	Count int         `json:"count"`
}

type RolesConfig struct {
	PlayersCount int                    `json:"playersCount"`
	RolesMp      map[string]*RoleConfig `json:"rolesMp"`
}

type ConfigsByPlayerCount []*RolesConfig

var (
	// FivePlayersConfigs represent configs with 5 players
	FivePlayersConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 5,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 5,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
			},
		},
	}

	// SixPlayersConfigs represent configs with 6 players
	SixPlayersConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 6,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 6,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 6,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
			},
		},
	}

	// SevenPlayersConfigs represent configs with 7 players
	SevenPlayersConfigs = &ConfigsByPlayerCount{
		// One active peaceful role
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		// Two active peaceful roles
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 3,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
	}

	// EightPlayersConfigs represent configs with 8 players
	EightPlayersConfigs = &ConfigsByPlayerCount{
		// Two active peaceful roles
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},

		// Three active peaceful roles
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
	}

	// NinePlayersConfigs represent configs with 9 players
	NinePlayersConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 9,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
	}

	// TenPlayersConfigs represent configs with 10 players
	TenPlayersConfigs = &ConfigsByPlayerCount{
		//Without maniac
		{
			PlayersCount: 10,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 2,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 10,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 3,
				},
			},
		},
		// With maniac
		{
			PlayersCount: 10,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  peaceful.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  peaceful.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  peaceful.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  peaceful.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  mafia.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  mafia.Don,
					Count: 1,
				},
				"Maniac": {
					Role:  maniac.Maniac,
					Count: 1,
				},
			},
		},
	}

	// ElevenPlayersConfigs represent configs with 11 players
	ElevenPlayersConfigs = &ConfigsByPlayerCount{}

	// TwelvePlayersConfigs represent configs with 12 players
	TwelvePlayersConfigs = &ConfigsByPlayerCount{}

	// ThirteenPlayersConfigs represent configs with 13 players
	ThirteenPlayersConfigs = &ConfigsByPlayerCount{}

	// FourteenPlayersConfigs represent configs with 14 players
	FourteenPlayersConfigs = &ConfigsByPlayerCount{}
)

var (
	// Configs int key represent count of players
	Configs = map[int]*ConfigsByPlayerCount{
		5:  FivePlayersConfigs,
		6:  SixPlayersConfigs,
		7:  SevenPlayersConfigs,
		8:  EightPlayersConfigs,
		9:  NinePlayersConfigs,
		10: TenPlayersConfigs,
		11: ElevenPlayersConfigs,
		12: TwelvePlayersConfigs,
		13: ThirteenPlayersConfigs,
		14: FourteenPlayersConfigs,
	}
)

func GetConfigsByPlayersCount(playersCount int) []*RolesConfig {
	return *Configs[playersCount]
}
