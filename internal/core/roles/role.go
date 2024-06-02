package roles

import "sort"

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

var MappedEmoji = map[string]string{
	"Citizen":   "",
	"Detective": "",
	"Doctor":    "",
	"Don":       "",
	"Fool":      "",
	"Mafia":     "",
	"Maniac":    "",
	"Peaceful":  "",
}

var StringTeam = map[Team]string{
	PeacefulTeam: ":black_heart: Peaceful",
	MafiaTeam:    ":heart: Mafia Team",
	ManiacTeam:   ":grey_heart: Maniac Team",
}

func GetEmojiByName(name string) string {
	return MappedEmoji[name]
}

func GetAllNightInteractionRolesNames() []string {
	var roles []string
	for name, role := range MappedRoles {
		if role.NightVoteOrder != -1 {
			roles = append(roles, name)
		}
	}
	return roles
}

func GetAllRolesNames() []string {
	var roles []string
	for name := range MappedRoles {
		roles = append(roles, name)
	}
	return roles
}

func GetRoleByName(roleName string) (*Role, bool) {
	role, ok := MappedRoles[roleName]
	return role, ok
}

func GetAllSortedRoles() []*Role {
	var allRoles []*Role
	for _, role := range MappedRoles {
		allRoles = append(allRoles, role)
	}

	sort.Slice(allRoles, func(i, j int) bool {
		return allRoles[i].Team < allRoles[j].Team
	})

	return allRoles
}

func GetStringTeam(team Team) string {
	return StringTeam[team]
}
