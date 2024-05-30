package maniac

import "github.com/https-whoyan/MafiaBot/internal/core/roles"

// Maniac description from Google.com
var (
	Maniac = &roles.Role{
		Name:           "Maniac",
		Team:           roles.ManiacTeam,
		NightVoteOrder: 3,
		Description: `
			The Maniac plays for himself. The Maniac's task is to get rid of 
			all the civilians and the Mafia. Every night he has the right to “kill” one player.`,
	}
)
