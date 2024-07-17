package roles

type Role struct {
	Name           string `json:"name" bson:"name"`
	Team           Team   `json:"team" bson:"team"`
	NightVoteOrder int    `json:"nightVoteOrder"`
	// Presents whether to execute immediately, the action of the role.
	UrgentCalculation bool
	// Allows for calculations to be made in the correct order after night.
	CalculationOrder int
	// Presents whether 2 player IDs are used in night actions of the role at once.
	IsTwoVotes  bool
	Description string `json:"description"`
}

var MappedRoles = map[string]*Role{
	"Citizen":   Citizen,
	"Detective": Detective,
	"Doctor":    Doctor,
	"Don":       Don,
	"Fool":      Fool,
	"Mafia":     Mafia,
	"Maniac":    Maniac,
	"Peaceful":  Peaceful,
	"Whore":     Whore,
}
