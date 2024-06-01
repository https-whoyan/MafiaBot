package roles

// Peaceful description from Google.com
var Peaceful = &Role{
	Name:           "Peaceful",
	Team:           PeacefulTeam,
	NightVoteOrder: -1,
	Description: `
		Peacekeeper is the most numerous role in the game. 
		Their job is to figure out the Mafia team players and eliminate them all on the day vote. 
		They don't go at night. They win when they eliminate all players not on their team.`,
}
