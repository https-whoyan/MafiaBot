package game

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
	myFMT "github.com/https-whoyan/MafiaBot/core/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// Night
// Actions with the game related to the night.
// Send will send signals to the channels about which role is currently voting. Comes from the g.run function
func (g *Game) Night(ch chan<- Signal) {
	select {
	case <-g.ctx.Done():
		return
	default:
		g.SetState(NightState)
		ch <- g.newSwitchStateSignal()
		g.Lock()
		g.NightCounter++
		g.Unlock()

		_, err := g.MainChannel.Write([]byte(g.getInitialNightMessage()))
		safeSendErrSignal(ch, err)

		// I'm getting the voting order
		g.RLock()
		orderToVote := g.RolesConfig.GetOrderToVote()
		g.RUnlock()

		// For each of the votes
		for _, votedRole := range orderToVote {
			g.RoleNightAction(votedRole, ch)
		}
		return
	}

}

/*
RoleNightAction
Counting variables, sending messages,
adding to spectators, and like that.
A follow-up call to the methods themselves is voice acceptance.
*/
func (g *Game) RoleNightAction(votedRole *rolesPack.Role, ch chan<- Signal) {
	select {
	case <-g.ctx.Done():
		return
	default:
		var err error

		g.Lock()
		g.NightVoting = votedRole
		g.Unlock()
		ch <- g.newSwitchVoteSignal()
		// Finding all the players with that role.
		// And finding interaction channel
		g.RLock()
		interactionChannel := g.RoleChannels[votedRole.Name]
		allPlayersWithRole := playerPack.SearchAllPlayersWithRole(g.Active, votedRole)
		g.RUnlock()

		voteDeadlineInt := myTime.VotingDeadline
		voteDeadline := time.Second * time.Duration(voteDeadlineInt)

		containsNotMutedPlayers := false

		// I go through each player and, with a mention, invite them to vote.
		// And if a player is locked, I tell him about it and add him to spectators for the duration of the vote.
		for _, voter := range allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				_, err = interactionChannel.Write([]byte(g.getMessageToPlayerThatIsMuted(voter)))
				safeSendErrSignal(ch, err)

				// Add to spectator
				err = channelPack.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(ch, err)

			} else {
				containsNotMutedPlayers = true
				_, err = interactionChannel.Write([]byte(g.getInvitingMessageToVote(voter, voteDeadlineInt)))
				safeSendErrSignal(ch, err)
			}
		}

		// From this differs in which channel the game will wait for the voice,
		//as well as the difference in the voice itself.
		switch !votedRole.IsTwoVotes {
		case true:
			g.oneVoteRoleNightVoting(allPlayersWithRole, containsNotMutedPlayers, voteDeadline, ch)
		default:
			g.twoVoterRoleNightVoting(allPlayersWithRole, containsNotMutedPlayers, voteDeadline, ch)
		}

		// Putting it back in the channel.
		for _, voter := range allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				err = channelPack.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(ch, err)
				_, err = interactionChannel.Write([]byte(g.getThanksToMutedPlayerMessage(voter)))
				safeSendErrSignal(ch, err)
			}
		}

		// Case when roles not need to urgent calculation
		// Return
		if !votedRole.UrgentCalculation {
			return
		}

		// Need to find a not empty vote.
		for _, voter := range allPlayersWithRole {
			voterVotesLen := len(voter.Votes)
			if voter.Votes[voterVotesLen-1] == EmptyVoteInt {
				continue
			}
			message := g.interaction(voter)
			if message != nil {
				_, err = interactionChannel.Write([]byte(*message))
				safeSendErrSignal(ch, err)
			}
			return
		}
		return
	}

}

/* The logic of accepting a role's vote, and timers, depending on whether the role votes with 2 votes or one. */

func (g *Game) oneVoteRoleNightVoting(allPlayersWithRole []*playerPack.Player,
	containsNotMutedPlayers bool, deadline time.Duration, ch chan<- Signal) {
	// Critic section with WaitGroup, don't use context completion check.
	var err error

	if !containsNotMutedPlayers {
		switch len(allPlayersWithRole) {
		case 0:
			ParralelierFullFakeVoteTimer(g.VoteChan)
			<-g.TwoVoteChan
		case 1:
			user := allPlayersWithRole[0]
			ParalleledFakeTimer(g.VoteChan, strconv.Itoa(user.ID), false)
			fakeVote := <-g.TwoVoteChan
			_ = g.NightTwoVote(fakeVote, nil)
		}
		return
	}

	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	for _, voter := range allPlayersWithRole {
		wg.Add(1)
		ParalleledVoteTimer(g.VoteChan, done, deadline,
			strconv.Itoa(voter.ID), false, wg)
	}
	for voteP := range g.VoteChan {
		err = g.NightOneVote(voteP, nil)
		if err == nil {
			break
		} else {
			ch <- newErrSignal(err)
		}
	}
	close(done)
	wg.Wait()
}

func (g *Game) twoVoterRoleNightVoting(allPlayersWithRole []*playerPack.Player,
	containsNotMutedPlayers bool, deadline time.Duration, ch chan<- Signal) {
	// Critic section with WaitGroup, don't use context completion check.
	var err error

	if !containsNotMutedPlayers {
		switch len(allPlayersWithRole) {
		case 0:
			ParralelierFullFakeTwoVotesTimer(g.TwoVoteChan)
			<-g.TwoVoteChan
		case 1:
			user := allPlayersWithRole[0]
			ParalleledTwoFakeTimer(g.TwoVoteChan, strconv.Itoa(user.ID), false)
			fakeVote := <-g.TwoVoteChan
			_ = g.NightTwoVote(fakeVote, nil)
		}
		return
	}

	// I create a channel for timers to work correctly.
	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	for _, voter := range allPlayersWithRole {
		wg.Add(1)
		ParalleledTwoVoteTimer(g.TwoVoteChan, done, deadline,
			strconv.Itoa(voter.ID), false, wg)
	}
	for voteP := range g.TwoVoteChan {
		err = g.NightTwoVote(voteP, nil)
		if err == nil {
			break
		} else {
			ch <- newErrSignal(err)
		}
	}
	close(done)
	wg.Wait()
}

// _________________
// Messages
// _________________

func (g *Game) getInitialNightMessage() string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr
	message := f.Bold("Night №") + f.Block(strconv.Itoa(g.NightCounter)) + " is coming." + f.LineSplitter()
	message += fmt.Sprintf("On this night you are played by %v players.", len(g.Active))
	return message
}

func (g *Game) getInvitingMessageToVote(p *playerPack.Player, deadlineInSeconds int) string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr
	message := f.Bold("Hello, " + f.Mention(p.ServerNick) + ". It's your turn to vote.")
	message += f.LineSplitter()
	message += myFMT.BoldUnderline(f, fmt.Sprintf("Deadline: %v seconds.", deadlineInSeconds))
	return message
}

func (g *Game) getMessageToPlayerThatIsMuted(p *playerPack.Player) string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr

	message := "Oops.... someone was muted today!" + f.Mention(p.ServerNick) +
		" , just chill, bro."
	return message
}

func (g *Game) getThanksToMutedPlayerMessage(p *playerPack.Player) string {
	g.RLock()
	defer g.RUnlock()
	message := g.fmtEr.Bold(g.fmtEr.Mention(p.ServerNick) + ", always thanks!")
	return message
}
