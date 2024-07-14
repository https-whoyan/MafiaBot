package config

import (
	"errors"
	"math"

	"github.com/https-whoyan/MafiaBot/core/roles"
)

// Presents all available role combinations and their number
// depending on the total number of players.
//
// Absolutely free for editing.

// Not to edit!
// Represent Config structs

type RoleConfig struct {
	Role  *roles.Role `json:"role" bson:"role"`
	Count int         `json:"count" bson:"count"`
}

type RolesConfig struct {
	PlayersCount int `json:"playersCount" bson:"playersCount"`
	// RolesMp present RoleConfig by RoleName.
	RolesMp map[string]*RoleConfig `json:"rolesMp" bson:"rolesMp"`
}

type ConfigsByPlayerCount []*RolesConfig

// Yes, edit it.

var (
	// TwoTestConfigs Test
	TwoTestConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 2,
			RolesMp: map[string]*RoleConfig{
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 2,
			RolesMp: map[string]*RoleConfig{
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 2,
			RolesMp: map[string]*RoleConfig{
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
			},
		},
	}

	// ThreeTestConfigs Test
	ThreeTestConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 3,
			RolesMp: map[string]*RoleConfig{
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 3,
			RolesMp: map[string]*RoleConfig{
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 2,
				},
			},
		},
		{
			PlayersCount: 3,
			RolesMp: map[string]*RoleConfig{
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
			},
		},
	}

	// FourTestConfigs Test
	FourTestConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 4,
			RolesMp: map[string]*RoleConfig{
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 1,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 4,
			RolesMp: map[string]*RoleConfig{
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 2,
				},
			},
		},
		{
			PlayersCount: 4,
			RolesMp: map[string]*RoleConfig{
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
			},
		},
	}

	// FivePlayersConfigs represent configs with 5 players
	FivePlayersConfigs = &ConfigsByPlayerCount{
		{
			PlayersCount: 5,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 5,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Mafia": {
					Role:  roles.Mafia,
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
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 6,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 6,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
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
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		// Two active peaceful roles
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 7,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 3,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
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
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},

		// Three active peaceful roles
		{
			PlayersCount: 8,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 3,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
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
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
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
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 10,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 3,
				},
			},
		},
		// With maniac
		{
			PlayersCount: 10,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 3,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
	}

	// ElevenPlayersConfigs represent configs with 11 players
	ElevenPlayersConfigs = &ConfigsByPlayerCount{
		// Without Maniac
		{
			PlayersCount: 11,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		{
			PlayersCount: 11,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		{
			PlayersCount: 11,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		// With Maniac
		{
			PlayersCount: 11,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
	}

	// TwelvePlayersConfigs represent configs with 12 players
	TwelvePlayersConfigs = &ConfigsByPlayerCount{
		// Same as eleven, but +1 peaceful
		// Without Maniac
		{
			PlayersCount: 12,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 6,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		{
			PlayersCount: 12,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		{
			PlayersCount: 12,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		// With Maniac
		{
			PlayersCount: 12,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
	}

	// ThirteenPlayersConfigs represent configs with 13 players
	ThirteenPlayersConfigs = &ConfigsByPlayerCount{
		// Without Maniac
		{
			PlayersCount: 13,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 6,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		// With Maniac
		// With Fool
		{
			PlayersCount: 13,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 13,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 4,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		// Without Fool
		{
			PlayersCount: 13,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 3,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 13,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
	}

	// FourteenPlayersConfigs represent configs with 14 players
	FourteenPlayersConfigs = &ConfigsByPlayerCount{
		// Same as 13
		// But +1 for peaceful
		// Without Maniac
		{
			PlayersCount: 14,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 7,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
			},
		},
		// With Maniac
		// With Fool
		{
			PlayersCount: 14,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 6,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 14,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 5,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Fool": {
					Role:  roles.Fool,
					Count: 1,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		// Without Fool
		{
			PlayersCount: 14,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 6,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 3,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
		{
			PlayersCount: 14,
			RolesMp: map[string]*RoleConfig{
				"Peaceful": {
					Role:  roles.Peaceful,
					Count: 6,
				},
				"Doctor": {
					Role:  roles.Doctor,
					Count: 1,
				},
				"Citizen": {
					Role:  roles.Citizen,
					Count: 1,
				},
				"Whore": {
					Role:  roles.Whore,
					Count: 1,
				},
				"Detective": {
					Role:  roles.Detective,
					Count: 1,
				},
				"Don": {
					Role:  roles.Don,
					Count: 1,
				},
				"Mafia": {
					Role:  roles.Mafia,
					Count: 2,
				},
				"Maniac": {
					Role:  roles.Maniac,
					Count: 1,
				},
			},
		},
	}
)

var (
	// Configs int key represent count of players
	Configs = map[int]*ConfigsByPlayerCount{
		2:  TwoTestConfigs,
		3:  ThreeTestConfigs,
		4:  FourTestConfigs,
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

// _______!!!!!!!!!!_______
// Not to edit!
// _______!!!!!!!!!!_______

var (
	SmallCountOfPeopleToConfig = errors.New("small people to config")
	BigCountOfPeopleToConfig   = errors.New("big people to config")
)

// GetConfigsByPlayersCount int in out present nearest available number of players
func GetConfigsByPlayersCount(playersCount int) ([]*RolesConfig, int, error) {
	smallestGroupOfPeople := GetMinPlayersCount()
	biggestGroupOfPeople := GetMaxPlayersCount()
	if playersCount < smallestGroupOfPeople {
		return nil, smallestGroupOfPeople, SmallCountOfPeopleToConfig
	} else if playersCount > GetMaxPlayersCount() {
		return nil, biggestGroupOfPeople, BigCountOfPeopleToConfig
	}
	return *Configs[playersCount], playersCount, nil
}

func GetMinPlayersCount() int {
	minPlayersCount := math.MaxInt
	for playersCount := range Configs {
		minPlayersCount = min(minPlayersCount, playersCount)
	}
	return minPlayersCount
}

func GetMaxPlayersCount() int {
	minPlayersCount := 0
	for playersCount := range Configs {
		minPlayersCount = max(minPlayersCount, playersCount)
	}
	return minPlayersCount
}
