package game

import "github.com/https-whoyan/MafiaBot/core/player"

// ResetTheVotes use to reset all player votes
func (g *Game) ResetTheVotes() {
	g.Lock()
	defer g.Unlock()
	allPlayers := g.Active

	for _, activePlayer := range allPlayers {
		activePlayer.DayVote = EmptyVoteInt
	}
}

// ResetAllInteractionsStatuses use to reset all player nightInteraction statuses
func (g *Game) ResetAllInteractionsStatuses() {
	g.Lock()
	defer g.Unlock()
	allPlayers := g.Active

	for _, activePlayer := range allPlayers {
		activePlayer.InteractionStatus = player.Passed
	}
}
