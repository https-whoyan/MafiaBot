package roles

// Maniac description from Google.com
var Maniac = &Role{
	Name:             "Maniac",
	Team:             ManiacTeam,
	NightVoteOrder:   3,
	CalculationOrder: 2,
	Description: `
		The Maniac plays for himself. The Maniac's task is to get rid of 
		all the civilians and the Mafia. Every night he has the right to “kill” one player.`,
}
