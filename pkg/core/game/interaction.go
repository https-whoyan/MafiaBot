package game

import (
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
	"strconv"
)

// All interactions for roles are declared here to avoid cyclic import.

type Message string

func (g *Game) interaction(p *player.Player) Message {
	if p.Role.NightVoteOrder == -1 {
		return ""
	}
	switch p.Role {
	case roles.Peaceful:
		return ""
	case roles.Fool:
		return ""
	case roles.Mafia:
		g.mafiaInteraction(p)
		return ""
	case roles.Doctor:
		g.doctorInteraction(p)
		return ""
	case roles.Don:
		return g.donInteraction(p)
	case roles.Detective:
		return g.detectiveInteraction(p)
	case roles.Whore:
		g.whoreInteraction(p)
		return ""
	case roles.Maniac:
		// Same action with mafia
		g.mafiaInteraction(p)
		return ""
	case roles.Citizen:
		g.citizenInteraction(p)
	}
	return ""
}

// ________________
// Mafia
// ________________

func (g *Game) mafiaInteraction(mafia *player.Player) {
	g.Lock()
	defer g.Unlock()
	nextDeadPlayer, isEmpty := g.interactionHelper(mafia)
	if isEmpty {
		return
	}

	nextDeadPlayer.LifeStatus = player.Dead
}

func (g *Game) doctorInteraction(doctor *player.Player) {
	g.Lock()
	defer g.Unlock()
	toVotedPlayer, isEmpty := g.interactionHelper(doctor)
	if isEmpty {
		return
	}

	if toVotedPlayer.LifeStatus == player.Dead {
		toVotedPlayer.LifeStatus = player.Alive
	}
}

func (g *Game) donInteraction(don *player.Player) Message {
	g.Lock()
	defer g.Unlock()
	f := g.fmtEr

	checkedPlayer, isEmpty := g.interactionHelper(don)
	if isEmpty {
		return ""
	}

	checkedPlayerRoleName := checkedPlayer.Role.Name

	message := "Checked player " + f.Block(strconv.Itoa(checkedPlayer.ID)) + ", role: " +
		g.fmtEr.Block(checkedPlayerRoleName)
	return Message(message)
}

func (g *Game) detectiveInteraction(detective *player.Player) Message {
	g.Lock()
	defer g.Unlock()

	checkedID1 := detective.Votes[len(detective.Votes)-1]
	checkedID2 := detective.Votes[len(detective.Votes)-2]

	if checkedID1 == EmptyVoteInt && checkedID2 == EmptyVoteInt {
		return ""
	}

	f := g.fmtEr
	checkedPlayer1 := player.SearchPlayerByGameID(g.Active, strconv.Itoa(checkedID1))
	checkedPlayer2 := player.SearchPlayerByGameID(g.Active, strconv.Itoa(checkedID2))

	isEqualsTeams := checkedPlayer1.Role.Team == checkedPlayer2.Role.Team

	var message string

	if isEqualsTeams {
		message = "Youu, players with id's " + f.Block(strconv.Itoa(checkedPlayer1.ID)) + ", " +
			f.Block(strconv.Itoa(checkedPlayer2.ID)) + f.Bold(" in one team.")
	} else {
		message = "Players with id's " + f.Block(strconv.Itoa(checkedPlayer1.ID)) + ", " +
			f.Block(strconv.Itoa(checkedPlayer2.ID)) + f.Bold(" in different team.")
	}
	return Message(message)
}

func (g *Game) whoreInteraction(whore *player.Player) {
	g.Lock()
	defer g.Unlock()
	mutedPlayer, isEmpty := g.interactionHelper(whore)
	if !isEmpty {
		return
	}
	mutedPlayer.InteractionStatus = player.Muted
}

func (g *Game) citizenInteraction(citizen *player.Player) {
	g.Lock()
	defer g.Unlock()
	defendedPlayer, isEmpty := g.interactionHelper(citizen)
	if !isEmpty {
		return
	}

	// Citizen is calculated by the most recent, then, if a civilian was killed, her
	// status would definitely be dead.
	defendedPlayer.LifeStatus = player.Alive
	if citizen.LifeStatus == player.Dead {
		defendedPlayer.LifeStatus = player.Dead
	}
}

// Helper
func (g *Game) interactionHelper(p *player.Player) (toVoted *player.Player, isEmpty bool) {
	lastVote := p.Votes[len(p.Votes)-1]

	if lastVote == EmptyVoteInt {
		isEmpty = true
		return
	}
	toVoted = player.SearchPlayerByGameID(g.Active, strconv.Itoa(lastVote))
	return
}
