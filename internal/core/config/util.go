package config

import (
	"math/rand"

	"github.com/https-whoyan/MafiaBot/internal/core/roles"
)

func GetShuffledRolesConfig(cfg *RolesConfig) []*roles.Role {
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
