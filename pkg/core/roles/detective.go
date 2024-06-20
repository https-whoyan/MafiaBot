package roles

// Detective description from my head)))
var Detective = &Role{
	Name:              "Detective",
	Team:              PeacefulTeam,
	UrgentCalculation: true,
	IsTwoVotes:        true,
	NightVoteOrder:    5,
	Description: `
		The commissioner checks 2 players at night, and finds out if they are 
		on the same team or not. Plays for peaceful players.`,
}
