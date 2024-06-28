package game

import (
	"fmt"
	"github.com/https-whoyan/MafiaBot/core/roles"
	"testing"

	"github.com/https-whoyan/MafiaBot/core/config"

	"github.com/stretchr/testify/assert"
)

func Test_Config_correct(t *testing.T) {
	for exceptedPlayersCount, configs := range config.Configs {
		for i, cfg := range *configs {

			// (Thanks Goland for name)
			t.Run(fmt.Sprintf("exceptedPlayersCount=%d, cfg=%+v", exceptedPlayersCount, cfg), func(t *testing.T) {
				assert.Equal(t, exceptedPlayersCount, cfg.PlayersCount,
					"error config playersCount, index: %d", i)

				// Check sum.
				sumOfPlayers := 0
				for _, roleConfiguration := range cfg.RolesMp {
					sumOfPlayers += roleConfiguration.Count
				}
				assert.Equal(t, exceptedPlayersCount, sumOfPlayers,
					"error config playersCount sum, index: %d", i)

				// Check for duplicate keys
				mpKeys := make(map[string]bool)
				for roleName, roleConfiguration := range cfg.RolesMp {
					if mpKeys[roleName] {
						// Another thanks to Goland)
						assert.Failf(t,
							fmt.Sprintf("duplicated role name %s in %+v; index: %+v",
								roleName, *roleConfiguration, i),
							"")
					}
					mpKeys[roleName] = true
				}

				// Check for duplicate roles
				mpContainsRole := make(map[*roles.Role]bool)
				for _, roleConfiguration := range cfg.RolesMp {
					checkedRole := roleConfiguration.Role
					if mpContainsRole[checkedRole] {
						assert.Failf(t, fmt.Sprintf("duplicated role, index: %d", i), "")
					}
				}

				defer func() {
					if r := recover(); r != nil {
						assert.Failf(t, fmt.Sprintf("recovered from panic in config at index: %d", i), "")
					}
				}()

				// Try to init game
				_, err := initHelper(cfg)
				if err != nil {
					t.Fatal(err)
				}
			})

		}
	}
}
