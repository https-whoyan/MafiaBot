package game

import (
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
)

// ClearDayVotes use to reset all player votes
func (g *Game) ClearDayVotes() {
	g.Lock()
	defer g.Unlock()
	allPlayers := g.Active

	for _, activePlayer := range *allPlayers {
		activePlayer.DayVote = EmptyVoteInt
	}
}

func (g *Game) ResetAllInteractionsStatuses() {
	g.Lock()
	defer g.Unlock()
	allPlayers := g.Active

	for _, activePlayer := range *allPlayers {
		activePlayer.InteractionStatus = player.Passed
	}
}

// UnderstandWinnerTeam is used to define scenarios where one team has 100% already won.
// A prime example is the 3/3 vote.
func (g *Game) UnderstandWinnerTeam() *roles.Team {
	g.RLock()
	defer g.RUnlock()
	allPlayers := g.Active

	// int represent count of players by their team
	teamsMp := make(map[roles.Team]int)
	for _, activePlayer := range *allPlayers {
		teamsMp[activePlayer.Role.Team]++
	}

	switch len(teamsMp) {
	case 1:
		var winnerTeam roles.Team
		for team := range teamsMp {
			winnerTeam = team
			break
		}
		return &winnerTeam
	case 2:
		// In this case, it all depends on the number of peaceful people.
		peacefulCnt, isNonZero := teamsMp[roles.PeacefulTeam]
		if !isNonZero { // (Is zero)
			// So that leaves the mafia team and the maniac team.
			// But the mafia team knows each other. So the mafia wins.
			winnerTeam := roles.MafiaTeam
			return &winnerTeam
		}
		// We find out the number of people from the opposite team.
		var (
			anotherTeam    roles.Team
			anotherTeamCnt int
		)
		for team, playersCnt := range teamsMp {
			if team != roles.PeacefulTeam {
				anotherTeamCnt = playersCnt
				anotherTeam = team
				break
			}
		}
		// Compared to the number of civilians
		if anotherTeamCnt >= peacefulCnt {
			return &anotherTeam
		}
		// return nil
		return nil
	}

	// 3 or more
	return nil
}
