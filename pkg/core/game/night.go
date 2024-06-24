package game

import (
	"fmt"
	"github.com/https-whoyan/MafiaBot/core/channel"
	myFMT "github.com/https-whoyan/MafiaBot/core/fmt"
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
	myTime "github.com/https-whoyan/MafiaBot/core/time"
	"strconv"
	"sync"
	"time"
)

// Night
// Actions with the game related to the night.
// Send will send signals to the channels about which role is currently voting. Comes from the g.run function
func (g *Game) Night(ch chan<- Signal) {
	select {
	case <-g.ctx.Done():
		return
	default:
		g.Lock()
		g.SetState(NightState)
		ch <- g.newSwitchStateSignal()
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
			g.roleNightAction(votedRole, ch)
		}
	}

}

/*
roleNightAction
Counting variables, sending messages,
adding to spectators, and like that.
A follow-up call to the methods themselves is voice acceptance.
*/
func (g *Game) roleNightAction(votedRole *roles.Role, ch chan<- Signal) {
	select {
	case <-g.ctx.Done():
		return
	default:
		var err error

		g.Lock()
		g.NightVoting = votedRole
		ch <- g.newSwitchVoteSignal()
		g.Unlock()
		// Finding all the players with that role.
		// And finding interaction channel
		g.RLock()
		interactionChannel := g.RoleChannels[votedRole.Name]
		allPlayersWithRole := player.SearchAllPlayersWithRole(g.Active, votedRole)
		g.RUnlock()

		voteDeadlineInt := myTime.VotingDeadline
		voteDeadline := time.Second * time.Duration(voteDeadlineInt)

		containsNotMutedPlayers := false

		// I go through each player and, with a mention, invite them to vote.
		// And if a player is locked, I tell him about it and add him to spectators for the duration of the vote.
		for _, voter := range allPlayersWithRole {
			if voter.InteractionStatus == player.Muted {
				_, err = interactionChannel.Write([]byte(g.getMessageToPlayerThatIsMuted(voter)))
				safeSendErrSignal(ch, err)

				// Add to spectator
				err = channel.FromUserToSpectator(interactionChannel, voter.Tag)
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
			if voter.InteractionStatus == player.Muted {
				err = channel.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(ch, err)
				_, err = interactionChannel.Write([]byte(g.getThanksToMutedPlayerMessage(voter)))
				safeSendErrSignal(ch, err)
			}
		}
	}

}

/* The logic of accepting a role's vote, and timers, depending on whether the role votes with 2 votes or one. */

func (g *Game) oneVoteRoleNightVoting(allPlayersWithRole []*player.Player,
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
		switch err {
		case nil:
			for i := 0; i <= len(allPlayersWithRole)-1; i++ {
				done <- struct{}{}
			}
			wg.Wait()
			close(done)
			break
		default:
			ch <- newErrSignal(err)
		}
	}
}

func (g *Game) twoVoterRoleNightVoting(allPlayersWithRole []*player.Player,
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
		switch err {
		case nil:
			for i := 0; i <= len(allPlayersWithRole)-1; i++ {
				done <- struct{}{}
			}
			wg.Wait()
			close(done)
			break
		default:
			ch <- newErrSignal(err)
		}
	}
}

// _________________
// Messages
// _________________

func (g *Game) getInitialNightMessage() string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr
	message := f.Bold("Night №") + f.Block(strconv.Itoa(g.NightCounter)) + " is coming." + f.LineSplitter()
	message += fmt.Sprintf("On this night you are played by %v players.", len(g.Active)) +
		f.LineSplitter() + f.LineSplitter()
	message += f.Italic("We wish you the best of luck)")
	return message
}

func (g *Game) getInvitingMessageToVote(p *player.Player, deadlineInSeconds int) string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr
	message := f.Bold("Hello, " + f.Mention(p.ServerNick) + ". It's your turn to vote.")
	message += f.LineSplitter()
	message += myFMT.BoldUnderline(f, fmt.Sprintf("Deadline: %v seconds.", deadlineInSeconds))
	return message
}

func (g *Game) getMessageToPlayerThatIsMuted(p *player.Player) string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr

	message := "Oops.... someone was muted today!" + f.Mention(p.ServerNick) +
		" , just chill, bro."
	return message
}

func (g *Game) getThanksToMutedPlayerMessage(p *player.Player) string {
	g.RLock()
	defer g.RUnlock()
	message := g.fmtEr.Bold(g.fmtEr.Mention(p.ServerNick) + ", always thanks!")
	return message
}
