package config

import (
	"math/rand"
	"sort"

	"github.com/https-whoyan/MafiaBot/core/converter"
	"github.com/https-whoyan/MafiaBot/core/roles"
)

func (cfg *RolesConfig) GetShuffledRolesConfig() []*roles.Role {
	var rolesArr []*roles.Role
	for _, roleConfig := range cfg.RolesMp {
		roleCount := roleConfig.Count
		role := roleConfig.Role
		for i := 1; i <= roleCount; i++ {
			rolesArr = append(rolesArr, role)
		}
	}

	rand.Shuffle(cfg.PlayersCount, func(i, j int) {
		rolesArr[i], rolesArr[j] = rolesArr[j], rolesArr[i]
	})

	return rolesArr
}

func (cfg *RolesConfig) GetTeamsByCfg() []roles.Team {
	mpTeams := make(map[roles.Team]bool)
	for _, role := range cfg.RolesMp {
		mpTeams[role.Role.Team] = true
	}

	teamsSlice := converter.GetMapKeys(mpTeams)
	sort.Slice(teamsSlice, func(i, j int) bool {
		return teamsSlice[i] < teamsSlice[j]
	})
	return teamsSlice
}

func (cfg *RolesConfig) GetMapKeyByTeamValuesRoleCfg() map[roles.Team][]*RoleConfig {
	teamsMp := make(map[roles.Team][]*RoleConfig)

	for _, roleCfg := range cfg.RolesMp {
		teamsMp[roleCfg.Role.Team] = append(teamsMp[roleCfg.Role.Team], roleCfg)
	}

	return teamsMp
}

func (cfg *RolesConfig) GetPlayersCountByTeam(team roles.Team) int {
	var count int
	for _, role := range cfg.RolesMp {
		if role.Role.Team == team {
			count += role.Count
		}
	}
	return count
}

func (cfg *RolesConfig) HasFool() bool {
	for _, roleCfg := range cfg.RolesMp {
		if roleCfg.Role == roles.Fool {
			return true
		}
	}
	return false
}
