package roles

// Doctor description from Google.com
var Doctor = &Role{
	Name:             "Doctor",
	Team:             PeacefulTeam,
	NightVoteOrder:   7,
	CalculationOrder: 3,
	Description: `
		The Doctor has the ability to heal the people of the town. 
		Each night, the Doctor tries to guess who was shot by the mafia and points
		that player to the host. If the Doctor guessed and “cured” the mafia victim, 
		the town wakes up without losses (or with fewer losses).`,
}
