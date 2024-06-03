package roles

import (
	"sort"
	"strings"
)

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
	"Whose":     Whore,
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
	"Whose":     "",
}

var StringTeam = map[Team]string{
	PeacefulTeam: ":heart: Peaceful",
	MafiaTeam:    ":black_heart: Mafia Team",
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
	allRoles := GetAllSortedRoles()
	var roleNames []string
	for _, role := range allRoles {
		roleNames = append(roleNames, role.Name)
	}

	return roleNames
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

func GetDefinitionOfRole(roleName string) string {
	fixDescription := func(s string) string {
		words := strings.Split(s, " ")
		return strings.Join(words, " ")
	}

	role := MappedRoles[roleName]
	name := "**__" + role.Name + "__**" + "\n"
	team := "**" + "Team: " + "**" + StringTeam[role.Team]
	description := fixDescription(role.Description)
	return name + "\n" + team + "\n" + description + "\n"
}
