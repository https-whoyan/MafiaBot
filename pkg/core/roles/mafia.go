package roles

// Mafia description from Google.com
var Mafia = &Role{
	Name:             "Mafia",
	Team:             MafiaTeam,
	NightVoteOrder:   5,
	CalculationOrder: 1,
	Description: `
		The goal of the mafia is to exterminate all civilians, or at least 
		stay with them in equal numbers. During the day the mafia 
		pretends to be honest townspeople, and at night the mafia cautiously wake up 
		and together choose a victim to “shoot”.`,
}
