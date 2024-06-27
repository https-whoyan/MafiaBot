package roles

var Citizen = &Role{
	Name:              "Citizen",
	Team:              PeacefulTeam,
	UrgentCalculation: true,
	CalculationOrder:  4,
	NightVoteOrder:    2,
	Description: `
		Hides a player with her at night, and the player is invulnerable to 
		the mafia and maniac that night, but if the citizen eventually 
		dies (she gets killed at night), the person she was hiding with her 
		will also die, so two people get killed.
			`,
}
