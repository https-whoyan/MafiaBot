package roles

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
	Description    string `json:"description"`
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
}

func GetRoleByName(roleName string) (*Role, bool) {
	role, ok := MappedRoles[roleName]
	return role, ok
}
