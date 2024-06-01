package players

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/https-whoyan/MafiaBot/internal/core/config"
	"github.com/https-whoyan/MafiaBot/internal/core/roles"
	"math/rand"
)

func generateRandomOrderToIDs(n int) []int {
	var IDs []int
	for i := 1; i <= n; i++ {
		IDs = append(IDs, i)
	}
	rand.Shuffle(n, func(i, j int) {
		IDs[i], IDs[j] = IDs[j], IDs[i]
	})

	return IDs
}

func getShuffledRolesConfig(cfg *config.RolesConfig) []*roles.Role {
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

func GeneratePlayers(users []*discordgo.User, cfg *config.RolesConfig) ([]*Player, error) {
	if len(users) != cfg.PlayersCount {
		return []*Player{}, errors.New("unexpected mismatch of playing participants and configs")
	}

	n := len(users)
	IDs := generateRandomOrderToIDs(n)
	rolesArr := getShuffledRolesConfig(cfg)

	players := make([]*Player, n)

	for i := 0; i <= n-1; i++ {
		players[i] = &Player{
			ID:                IDs[i],
			OldNick:           users[i].Username,
			Tag:               users[i],
			Role:              rolesArr[i],
			Vote:              -1,
			LifeStatus:        Alive,
			InteractionStatus: Passed,
		}
	}

	return players, nil
}
