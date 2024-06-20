package roles

// Whore description from Google.com
var Whore = &Role{
	Name:              "Whore",
	Team:              PeacefulTeam,
	NightVoteOrder:    1,
	UrgentCalculation: true,
	Description: `
		The role of the Prostitute is to choose one of the players to spend 
		the night with. The Prostitute blocks the actions of the chosen character - 
		the mafia doesn't shoot, the maniac doesn't kill, the doctor doesn't heal, the 
		sheriff doesn't check, and so on.”.`,
}
