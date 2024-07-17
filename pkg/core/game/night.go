package game

import (
	"sort"
	"time"

	channelPack "github.com/https-whoyan/MafiaBot/core/channel"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	rolesPack "github.com/https-whoyan/MafiaBot/core/roles"
	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// Night
// Actions with the game related to the night.
// Send will send signals to the channels about which role is currently voting. Comes from the g.run function
func (g *Game) Night() NightLog {
	select {
	case <-g.ctx.Done():
		return NightLog{}
	default:
		g.SetState(NightState)
		g.infoSender <- g.newSwitchStateSignal()

		err := g.Messenger.Night.SendInitialNightMessage(g.MainChannel)
		safeSendErrSignal(g.infoSender, err)

		// I'm getting the voting order
		g.RLock()
		orderToVote := g.RolesConfig.GetOrderToVote()
		g.RUnlock()

		// For each of the votes
		for _, votedRole := range orderToVote {
			// To avoid for shorting
			votedRoleClone := votedRole
			g.RoleNightAction(votedRoleClone)
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
			if p.Role.CalculationOrder > 0 && !p.Role.UrgentCalculation {
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
		return g.NewNightLog()
	}

}

/*
RoleNightAction

	Counting variables, sending messages,
	adding to spectators, and like that.

	A follow-up call to the methods themselves is voice acceptance.
*/
func (g *Game) RoleNightAction(votedRole *rolesPack.Role) {

	select {
	case <-g.ctx.Done():
		return
	default:
		var err error

		g.Lock()
		g.NightVoting = votedRole
		g.Unlock()
		g.infoSender <- g.newSwitchVoteSignal()
		// Finding all the players with that role.
		// And finding nightInteraction channel
		g.RLock()
		interactionChannel := g.RoleChannels[votedRole.Name]
		allPlayersWithRole := g.Active.SearchAllPlayersWithRole(votedRole)
		g.RUnlock()

		var (
			nonEmptyVote1 = EmptyVoteInt
			nonEmptyVote2 = EmptyVoteInt
		)

		sendToOtherEmptyVotes := func(nonEmptyVoter *playerPack.Player) {
			voterLen := len(nonEmptyVoter.Votes)
			if votedRole.IsTwoVotes {
				nonEmptyVote1 = nonEmptyVoter.Votes[voterLen-2]
				nonEmptyVote2 = nonEmptyVoter.Votes[voterLen-1]
			} else {
				nonEmptyVote1 = nonEmptyVoter.Votes[voterLen-1]
			}

			// Set to other empty votes
			for _, playerWithRole := range *allPlayersWithRole {
				if playerWithRole == nonEmptyVoter {
					continue
				}
				if votedRole.IsTwoVotes {
					playerWithRole.Votes = append(playerWithRole.Votes, nonEmptyVote1, nonEmptyVote2)
				} else {
					playerWithRole.Votes = append(playerWithRole.Votes, nonEmptyVote1)
				}
			}
			return
		}
		findOrStandNotEmptyVoter := func() (nonEmptyVoter *playerPack.Player) {
			// Need to find a not empty Vote.
			for _, voter := range *allPlayersWithRole {
				voterVotesLen := len(voter.Votes)
				if voterVotesLen == 0 {
					continue
				}
				if voter.Votes[voterVotesLen-1] == EmptyVoteInt {
					continue
				}
				nonEmptyVoter = voter
				break
			}

			// If deadline pass, or player set EmptyVote, stand nonEmptyVoter
			// on somebody, and we'll put votes on the player.
			if nonEmptyVoter == nil {
				for _, p := range *allPlayersWithRole {
					nonEmptyVoter = p

					if votedRole.IsTwoVotes {
						p.Votes = append(p.Votes, nonEmptyVote1, nonEmptyVote2)
					} else {
						p.Votes = append(p.Votes, nonEmptyVote1)
					}

					break
				}
			}
			return
		}

		voteDeadlineInt := myTime.VotingDeadline
		voteDeadline := time.Second * time.Duration(voteDeadlineInt)

		containsNotMutedPlayers := false

		// I go through each player and, with a mention, invite them to Vote.
		// And if a player is locked, I tell him about it and add him to spectators for the duration of the Vote.
		for _, voter := range *allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				err = g.Messenger.Night.SendToPlayerThatIsMutedMessage(voter, interactionChannel)
				safeSendErrSignal(g.infoSender, err)

				// Add to spectator
				err = channelPack.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(g.infoSender, err)

			} else {
				containsNotMutedPlayers = true
				err = g.Messenger.Night.SendInvitingToVoteMessage(voter, voteDeadlineInt, interactionChannel)
				safeSendErrSignal(g.infoSender, err)
			}
		}

		// From this differs in which channel the game will wait for the voice,
		//as well as the difference in the voice itself.
		switch votedRole.IsTwoVotes {
		case true:
			g.twoVoterRoleNightVoting(containsNotMutedPlayers, voteDeadline)
		case false:
			g.oneVoteRoleNightVoting(containsNotMutedPlayers, voteDeadline)
		}

		// Putting it back in the channel.
		for _, voter := range *allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				err = channelPack.FromSpectatorToUser(interactionChannel, voter.Tag)
				safeSendErrSignal(g.infoSender, err)

				err = g.Messenger.Night.SendThanksToMutedPlayerMessage(voter, interactionChannel)
				safeSendErrSignal(g.infoSender, err)
			}
		}

		nonEmptyVoter := findOrStandNotEmptyVoter()
		sendToOtherEmptyVotes(nonEmptyVoter)

		// Case when roles need to urgent calculation
		if votedRole.UrgentCalculation {
			message := g.nightInteraction(nonEmptyVoter)
			if message != nil {
				_, err = interactionChannel.Write([]byte(*message))
				safeSendErrSignal(g.infoSender, err)
			}
		}
	}
}

/*
	The logic of accepting a role's Vote, and timers,
	depending on whether the role votes with 2 votes or one.
*/

func (g *Game) waitOneVoteRoleFakeTimer() {
	g.randomTimer()

	for {
		isNeedToContinue := false
		select {
		case voteP := <-g.VoteChan:
			// All votes will be with errors
			err := g.NightOneVote(voteP, nil)
			g.infoSender <- newErrSignal(err)
			isNeedToContinue = true
			break
		case <-g.timerDone:
			break
		case <-g.ctx.Done():
			break
		}

		if !isNeedToContinue {
			break
		}
	}
}

func (g *Game) oneVoteRoleNightVoting(containsNotMutedPlayers bool, deadline time.Duration) {
	var err error

	if !containsNotMutedPlayers {
		g.waitOneVoteRoleFakeTimer()
		return
	}

	g.timer(deadline)

	for {
		isNeedToContinue := false
		select {
		case voteP := <-g.VoteChan:
			err = g.NightOneVote(voteP, nil)
			if err == nil {
				g.timerStop <- struct{}{}
				break
			} else {
				g.infoSender <- newErrSignal(err)
				isNeedToContinue = true
				break
			}
		case <-g.timerDone:
			break
		case <-g.ctx.Done():
			break
		}
		if !isNeedToContinue {
			break
		}
	}
}

func (g *Game) waitTwoVoteRoleFakeTimer() {
	g.randomTimer()

	for {
		isNeedToContinue := false
		select {
		case voteP := <-g.TwoVoteChan:
			// All votes will be with errors
			err := g.NightTwoVote(voteP, nil)
			g.infoSender <- newErrSignal(err)
			isNeedToContinue = true
			break
		case <-g.timerDone:
			break
		case <-g.ctx.Done():
			break
		}

		if !isNeedToContinue {
			break
		}
	}
}

func (g *Game) twoVoterRoleNightVoting(containsNotMutedPlayers bool, deadline time.Duration) {
	var err error

	if !containsNotMutedPlayers {
		g.waitTwoVoteRoleFakeTimer()
		return
	}

	g.timer(deadline)

	for {
		isNeedToContinue := false
		select {
		case voteP := <-g.TwoVoteChan:
			err = g.NightTwoVote(voteP, nil)
			if err == nil {
				g.timerStop <- struct{}{}
				break
			} else {
				g.infoSender <- newErrSignal(err)
				isNeedToContinue = true
			}
		case <-g.timerDone:
			break
		case <-g.ctx.Done():
			break
		}

		if !isNeedToContinue {
			break
		}
	}
}

// AffectNight changes players according to the night's actions.
// Errors during execution are sent to the channel
func (g *Game) AffectNight(l NightLog) {
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
		var newDeadPersons []*playerPack.DeadPlayer

		for _, deadID := range l.Dead {
			g.Active.ToDead(playerPack.IDType(deadID), playerPack.KilledAtNight, g.NightCounter, g.Dead)
		}

		// I will add add add all killed players after a minute of players after a minute of
		// players after a minute, so, using goroutine.
		go func(newDeadPersons []*playerPack.DeadPlayer) {
			duration := myTime.LastWordDeadline * time.Second
			ticker := time.NewTicker(duration)

			g.RLock()
			mainChannel := g.MainChannel
			roleChannels := g.RoleChannels
			g.RUnlock()

			defer ticker.Stop()
			select {
			case <-g.ctx.Done():
				return
			case <-ticker.C:
				// I'm adding new dead players to the spectators in the channels (so they won't be so bored)
				for _, p := range newDeadPersons {
					for _, interactionChannel := range roleChannels {
						select {
						case <-g.ctx.Done():
							return
						default:
							safeSendErrSignal(g.infoSender, channelPack.FromUserToSpectator(interactionChannel, p.Tag))
							break
						}
					}
					select {
					case <-g.ctx.Done():
						return
					default:
						safeSendErrSignal(g.infoSender, channelPack.FromUserToSpectator(mainChannel, p.Tag))
						break
					}
				}
			}
		}(newDeadPersons)

		// Sending a message about who died today.
		err := g.Messenger.AfterNight.SendAfterNightMessage(l, g.MainChannel)
		safeSendErrSignal(g.infoSender, err)
		// Then, for each person try to do his reincarnation
		g.Unlock()
		for _, p := range *g.Active {
			g.reincarnation(p)
		}
		return
	}
}
