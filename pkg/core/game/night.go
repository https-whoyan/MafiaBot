package game

import (
	"github.com/https-whoyan/MafiaBot/core/converter"
	"sort"
	"strconv"
	"sync"
	"time"

	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
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

		err := g.messenger.Night.SendInitialNightMessage(g.MainChannel)
		safeSendErrSignal(ch, err)

		// I'm getting the voting order
		g.RLock()
		orderToVote := g.RolesConfig.GetOrderToVote()
		g.RUnlock()

		// For each of the votes
		for _, votedRole := range orderToVote {
			g.RoleNightAction(votedRole, ch)
		}
		// On this line, all votes are accepted.
		// I hereby signify that the voting is over.
		g.Lock()
		g.NightVoting = nil
		g.Unlock()
		g.RLock()
		// I do the rest of the interactions that come after the vote.
		var needToProcessPlayers []*playerPack.Player
		for _, p := range *g.Active {
			if p.Role.CalculationOrder > 0 {
				needToProcessPlayers = append(needToProcessPlayers, p)
			}
		}
		g.RUnlock()
		sort.Slice(needToProcessPlayers, func(i, j int) bool {
			return needToProcessPlayers[i].Role.CalculationOrder < needToProcessPlayers[j].Role.CalculationOrder
		})
		for _, p := range needToProcessPlayers {
			g.nightInteraction(p)
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
		// And finding nightInteraction channel
		g.RLock()
		interactionChannel := g.RoleChannels[votedRole.Name]
		allPlayersWithRole := g.Active.SearchAllPlayersWithRole(votedRole)
		g.RUnlock()

		voteDeadlineInt := myTime.VotingDeadline
		voteDeadline := time.Second * time.Duration(voteDeadlineInt)

		containsNotMutedPlayers := false

		// I go through each player and, with a mention, invite them to Vote.
		// And if a player is locked, I tell him about it and add him to spectators for the duration of the Vote.
		for _, voter := range *allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				err = g.messenger.Night.SendToPlayerThatIsMutedMessage(voter, interactionChannel)
				safeSendErrSignal(ch, err)

				// Add to spectator
				err = channelPack.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(ch, err)

			} else {
				containsNotMutedPlayers = true
				err = g.messenger.Night.SendInvitingToVoteMessage(voter, voteDeadlineInt, interactionChannel)
				safeSendErrSignal(ch, err)
			}
		}

		// From this differs in which channel the game will wait for the voice,
		//as well as the difference in the voice itself.
		slicePlayers := converter.GetMapValues(*allPlayersWithRole)
		switch votedRole.IsTwoVotes {
		case true:
			g.twoVoterRoleNightVoting(slicePlayers, containsNotMutedPlayers, voteDeadline, ch)
		case false:
			g.oneVoteRoleNightVoting(slicePlayers, containsNotMutedPlayers, voteDeadline, ch)
		}

		// Putting it back in the channel.
		for _, voter := range *allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				err = channelPack.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(ch, err)

				err = g.messenger.Night.SendThanksToMutedPlayerMessage(voter, interactionChannel)
				safeSendErrSignal(ch, err)
			}
		}

		// Case when roles not need to urgent calculation
		// Return
		if !votedRole.UrgentCalculation {
			return
		}

		// Need to find a not empty Vote.
		for _, voter := range *allPlayersWithRole {
			voterVotesLen := len(voter.Votes)
			if voter.Votes[voterVotesLen-1] == EmptyVoteInt {
				continue
			}
			message := g.nightInteraction(voter)
			if message != nil {
				_, err = interactionChannel.Write([]byte(*message))
				safeSendErrSignal(ch, err)
			}
			return
		}
		return
	}

}

/* The logic of accepting a role's Vote, and timers, depending on whether the role votes with 2 votes or one. */

func (g *Game) waitOneVoteRoleFakeTimer(allPlayersWithRole []*playerPack.Player, ch chan<- Signal) {
	done := make(chan struct{})
	FakeTimer(done)
	// All votes will be with errors.
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-done
		close(done)
	}()
	go func() {
		defer wg.Done()
		select {
		case <-done:
			return
		case voteP := <-g.VoteChan:
			select {
			case ch <- newErrSignal(g.NightOneVote(voteP, nil)):
			case <-done:
				return
			}
		}
	}()
	wg.Wait()
	// Set to all Empty votes
	for _, p := range allPlayersWithRole {
		p.Votes = append(p.Votes, EmptyVoteInt)
	}
}
func (g *Game) oneVoteRoleNightVoting(allPlayersWithRole []*playerPack.Player,
	containsNotMutedPlayers bool, deadline time.Duration, ch chan<- Signal) {
	// Critic section with WaitGroup, don't use context completion check.
	var err error

	if !containsNotMutedPlayers {
		g.waitOneVoteRoleFakeTimer(allPlayersWithRole, ch)
		return
	}

	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	for _, voter := range allPlayersWithRole {
		VoteTimer(g.VoteChan, done, deadline,
			strconv.Itoa(int(voter.ID)), false, wg)
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

func (g *Game) waitTwoVoteRoleFakeTimer(allPlayersWithRole []*playerPack.Player, ch chan<- Signal) {
	done := make(chan struct{})
	FakeTimer(done)
	// All votes will be with errors.
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-done
		close(done)
	}()
	go func() {
		defer wg.Done()
		select {
		case <-done:
			return
		case voteP := <-g.TwoVoteChan:
			select {
			case ch <- newErrSignal(g.NightTwoVote(voteP, nil)):
			case <-done:
				return
			}
		}
	}()
	wg.Wait()
	// Set to all Empty votes
	for _, p := range allPlayersWithRole {
		p.Votes = append(p.Votes, EmptyVoteInt)
	}
}
func (g *Game) twoVoterRoleNightVoting(allPlayersWithRole []*playerPack.Player,
	containsNotMutedPlayers bool, deadline time.Duration, ch chan<- Signal) {
	// Critic section with WaitGroup, don't use context completion check.
	var err error

	if !containsNotMutedPlayers {
		done := make(chan struct{})
		FakeTimer(done)
		// All votes will be with errors.
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-done
			close(done)
		}()
		go func() {
			defer wg.Done()
			select {
			case <-done:
				return
			case voteP := <-g.TwoVoteChan:
				select {
				case ch <- newErrSignal(g.NightTwoVote(voteP, nil)):
				case <-done:
					return
				}
			}
		}()
		wg.Wait()
		// Set to all Empty votes
		for _, p := range allPlayersWithRole {
			p.Votes = append(p.Votes, EmptyVoteInt, EmptyVoteInt)
		}
		return
	}

	// I create a channel for timers to work correctly.
	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	for _, voter := range allPlayersWithRole {
		TwoVoteTimer(g.TwoVoteChan, done, deadline,
			strconv.Itoa(int(voter.ID)), false, wg)
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

// AffectNight changes players according to the night's actions.
// Errors during execution are sent to the channel
func (g *Game) AffectNight(l NightLog, ch chan<- Signal) {
	// Clearing all statuses
	if !g.IsRunning() {
		panic("Game is not running")
	}
	if g.ctx == nil {
		panic("Game context is nil, then, don't initialed")
	}
	select {
	case <-g.ctx.Done():
		return
	default:
		g.ResetAllInteractionsStatuses()
		g.Lock()

		// Splitting arrays.
		newActivePlayers := make(playerPack.Players)
		var newDeadPersons []*playerPack.DeadPlayer

		for _, p := range *g.Active {
			if p.LifeStatus == playerPack.Dead {
				newDeadPlayer := playerPack.NewDeadPlayer(p, playerPack.KilledAtNight, g.NightCounter)
				newDeadPersons = append(newDeadPersons, newDeadPlayer)
			} else {
				newActivePlayers[p.ID] = p
			}
		}

		// I will add add add all killed players after a minute of players after a minute of
		// players after a minute, so, using goroutine.
		go func(newDeadPersons []*playerPack.DeadPlayer) {
			duration := myTime.LastWordDeadline * time.Second
			ticker := time.NewTicker(duration)
			defer ticker.Stop()
			select {
			case <-g.ctx.Done():
				return
			case <-ticker.C:
				g.RLock()
				defer g.RUnlock()
				// I'm adding new dead players to the spectators in the channels (so they won't be so bored)
				for _, p := range newDeadPersons {
					for _, interactionChannel := range g.RoleChannels {
						safeSendErrSignal(ch, channelPack.FromUserToSpectator(interactionChannel, p.Tag))
					}
					safeSendErrSignal(ch, channelPack.FromUserToSpectator(g.MainChannel, p.Tag))
				}
			}
		}(newDeadPersons)

		// Changing arrays according to the night
		g.Active = &newActivePlayers
		g.Dead.Add(newDeadPersons...)
		g.Unlock()

		// Sending a message about who died today.
		err := g.messenger.AfterNight.SendAfterNightMessage(l, g.MainChannel)
		safeSendErrSignal(ch, err)
		// Then, for each person try to do his reincarnation
		for _, p := range newActivePlayers {
			g.reincarnation(ch, p)
		}
		return
	}
}
