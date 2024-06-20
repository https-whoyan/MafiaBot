package roles

// Don description from Google.com
var Don = &Role{
	Name:              "Don",
	Team:              MafiaTeam,
	UrgentCalculation: true,
	NightVoteOrder:    4,
	Description: `
		The mafia don is the main mafioso, playing against honest citizens. 
		The don's role is almost identical to the mafia role, except that the don
		can check any player at night - and find out from the host
		if he is an active role.`,
}
