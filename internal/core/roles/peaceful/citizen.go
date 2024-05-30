package peaceful

import "github.com/https-whoyan/MafiaBot/internal/core/roles"

var (
	Citizen = &roles.Role{
		Name:           "Citizen",
		Team:           roles.PeacefulTeam,
		NightVoteOrder: 2,
		Description: `
			Hides a player with her at night, and the player is invulnerable to 
			the mafia and maniac that night, but if the citizen eventually 
			dies (she gets killed at night), the person she was hiding with her 
			will also die, so two people get killed.
			`,
	}
)
