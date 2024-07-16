package game

import (
	"strconv"

	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
)

// All interactions for roles are declared here to avoid cyclic import.

type InteractionMessage string

func (g *Game) nightInteraction(p *player.Player) *InteractionMessage {
	if p.Role.NightVoteOrder == -1 {
		return nil
	}
	switch p.Role {
	case roles.Peaceful:
		return nil
	case roles.Fool:
		return nil
	case roles.Mafia:
		g.mafiaInteraction(p)
		return nil
	case roles.Doctor:
		g.doctorInteraction(p)
		return nil
	case roles.Don:
		return g.donInteraction(p)
	case roles.Detective:
		return g.detectiveInteraction(p)
	case roles.Whore:
		g.whoreInteraction(p)
		return nil
	case roles.Maniac:
		g.maniacInteraction(p)
		return nil
	case roles.Citizen:
		g.citizenInteraction(p)
	}
	return nil
}

/* Mafia */

func (g *Game) mafiaInteraction(mafia *player.Player) {
	g.Lock()
	defer g.Unlock()
	nextDeadPlayer, isEmpty := g.interactionHelper(mafia)
	if isEmpty {
		return
	}

	nextDeadPlayer.LifeStatus = player.Dead
}
func (g *Game) donInteraction(don *player.Player) *InteractionMessage {
	g.Lock()
	defer g.Unlock()
	f := g.Messenger.f

	checkedPlayer, isEmpty := g.interactionHelper(don)
	if isEmpty {
		return nil
	}

	checkedPlayerRoleName := checkedPlayer.Role.Name

	message := InteractionMessage("Checked player " + f.Block(strconv.Itoa(int(checkedPlayer.ID))) + ", role: " +
		g.Messenger.f.Block(checkedPlayerRoleName))
	return &message
}

/* Maniac */

func (g *Game) maniacInteraction(maniac *player.Player) {
	// Same as mafia
	g.Lock()
	defer g.Unlock()
	nextDeadPlayer, isEmpty := g.interactionHelper(maniac)
	if isEmpty {
		return
	}

	nextDeadPlayer.LifeStatus = player.Dead
}

/* Peaceful */

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
func (g *Game) detectiveInteraction(detective *player.Player) *InteractionMessage {
	g.Lock()
	defer g.Unlock()

	checkedID1 := detective.Votes[len(detective.Votes)-1]
	checkedID2 := detective.Votes[len(detective.Votes)-2]

	if checkedID1 == EmptyVoteInt && checkedID2 == EmptyVoteInt {
		return nil
	}

	f := g.Messenger.f
	checkedPlayer1 := g.Active.SearchPlayerByGameID(strconv.Itoa(checkedID1))
	checkedPlayer2 := g.Active.SearchPlayerByGameID(strconv.Itoa(checkedID2))

	isEqualsTeams := checkedPlayer1.Role.Team == checkedPlayer2.Role.Team

	var message string

	if isEqualsTeams {
		message = "Players with id's " + f.Block(strconv.Itoa(int(checkedPlayer1.ID))) + ", " +
			f.Block(strconv.Itoa(int(checkedPlayer2.ID))) + f.Bold(" in one team.")
	} else {
		message = "Players with id's " + f.Block(strconv.Itoa(int(checkedPlayer1.ID))) + ", " +
			f.Block(strconv.Itoa(int(checkedPlayer2.ID))) + f.Bold(" in different team.")
	}
	typedMessage := InteractionMessage(message)
	return &typedMessage
}
func (g *Game) whoreInteraction(whore *player.Player) {
	g.Lock()
	defer g.Unlock()
	mutedPlayer, isEmpty := g.interactionHelper(whore)
	if isEmpty {
		return
	}
	mutedPlayer.InteractionStatus = player.Muted
}
func (g *Game) citizenInteraction(citizen *player.Player) {
	g.Lock()
	defer g.Unlock()
	defendedPlayer, isEmpty := g.interactionHelper(citizen)
	if isEmpty {
		return
	}

	// Citizen is calculated by the most recent, then, if a civilian was killed, her
	// status would definitely be dead.
	defendedPlayer.LifeStatus = player.Alive
	if citizen.LifeStatus == player.Dead {
		defendedPlayer.LifeStatus = player.Dead
	}
}

/* Helper */

func (g *Game) interactionHelper(p *player.Player) (toVoted *player.Player, isEmpty bool) {
	lastVote := p.Votes[len(p.Votes)-1]

	if lastVote == EmptyVoteInt {
		isEmpty = true
		return
	}
	toVoted = (*g.Active)[player.IDType(lastVote)]
	return
}
