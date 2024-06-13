package roles

import (
	"sort"
	"strings"

	"github.com/https-whoyan/MafiaBot/core/message/fmt"
)

// Utils.

// __________________
// This contains all the functions that link the role name to the role, and things like that.
// __________________

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

// _____________________________________________________________________
// Beautiful presentations of roles to display information about them.
// _____________________________________________________________________

func GetDefinitionOfRole(f fmt.FmtInterface, roleName string) string {
	fixDescription := func(s string) string {
		words := strings.Split(s, " ")
		return strings.Join(words, " ")
	}

	role := MappedRoles[roleName]
	var message string

	name := string(fmt.BoldUnderline(f, role.Name))
	team := string(f.Bold("Team: ")) + StringTeam[role.Team]
	description := fixDescription(role.Description)
	message = name + "\n\n" + team + "\n" + description
	return message
}
