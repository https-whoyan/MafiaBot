package peaceful

import "github.com/https-whoyan/MafiaBot/internal/core/roles"

// Detective description from my head)))
var (
	Detective = &roles.Role{
		Name:           "Detective",
		Team:           roles.PeacefulTeam,
		NightVoteOrder: 5,
		Description: `
			The commissioner checks 2 players at night, and finds out if they are 
			on the same team or not. Plays for peaceful players.
			`,
	}
)
