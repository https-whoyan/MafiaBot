package roles

// All information regarding roles.

type Team int

const (
	PeacefulTeam Team = 1
	MafiaTeam    Team = 2
	ManiacTeam   Team = 3
)

type Role struct {
	Name           string `json:"name"`
	Team           Team   `json:"team"`
	NightVoteOrder int    `json:"nightVoteOrder"`
	// Presents whether to execute immediately, the action of the role.
	UrgentCalculation bool
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
