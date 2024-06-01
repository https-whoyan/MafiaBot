package roles

// Fool description from Google.com
var Fool = &Role{
	Name:           "Fool",
	Team:           PeacefulTeam,
	NightVoteOrder: -1,
	Description: `
		The Fool plays by himself. 
		There are no night game actions for this role.
		The Fool must convince the townspeople to execute him 
		(for example, pretending to be a mobster or a maniac). This is the only way he can win. 
		After he is executed, the game ends - all other participants lose. If the Fool is killed at night, he loses.
		The role of the Fool is available in automatic games without bots, which take place every hour.`,
}
