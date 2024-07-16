package game

import (
	"github.com/https-whoyan/MafiaBot/core/converter"
	"sort"
	"strconv"
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
			g.RoleNightAction(votedRole)
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
		slicePlayers := converter.GetMapValues(*allPlayersWithRole)
		switch votedRole.IsTwoVotes {
		case true:
			g.twoVoterRoleNightVoting(slicePlayers, containsNotMutedPlayers, voteDeadline)
		case false:
			g.oneVoteRoleNightVoting(slicePlayers, containsNotMutedPlayers, voteDeadline)
		}

		// Putting it back in the channel.
		for _, voter := range *allPlayersWithRole {
			if voter.InteractionStatus == playerPack.Muted {
				err = channelPack.FromUserToSpectator(interactionChannel, voter.Tag)
				safeSendErrSignal(g.infoSender, err)

				err = g.Messenger.Night.SendThanksToMutedPlayerMessage(voter, interactionChannel)
				safeSendErrSignal(g.infoSender, err)
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
				safeSendErrSignal(g.infoSender, err)
			}
			return
		}
		return
	}

}

/* The logic of accepting a role's Vote, and timers, depending on whether the role votes with 2 votes or one. */

func (g *Game) waitOneVoteRoleFakeTimer(allPlayersWithRole []*playerPack.Player) {
	g.randomTimer()

	select {
	case voteP := <-g.VoteChan:
		// All votes will be with errors
		err := g.NightOneVote(voteP, nil)
		g.infoSender <- newErrSignal(err)
	case <-g.timerDone:
		votedPlayer := int(allPlayersWithRole[0].ID)
		voteProvider := NewVoteProvider(strconv.Itoa(votedPlayer), EmptyVoteStr, false)
		_ = g.NightOneVote(voteProvider, nil)
		break
	}
	for _, p := range allPlayersWithRole {
		p.Votes = append(p.Votes, EmptyVoteInt)
	}
}

func (g *Game) oneVoteRoleNightVoting(allPlayersWithRole []*playerPack.Player,
	containsNotMutedPlayers bool, deadline time.Duration) {
	// Critic section with WaitGroup, don't use context completion check.
	var err error

	if !containsNotMutedPlayers {
		g.waitOneVoteRoleFakeTimer(allPlayersWithRole)
		return
	}

	g.timer(deadline)

	select {
	case voteP := <-g.VoteChan:
		err = g.NightOneVote(voteP, nil)
		if err == nil {
			g.timerStop <- struct{}{}
			break
		} else {
			g.infoSender <- newErrSignal(err)
		}
	case <-g.timerDone:
		votedPlayer := int(allPlayersWithRole[0].ID)
		voteProvider := NewVoteProvider(strconv.Itoa(votedPlayer), EmptyVoteStr, false)
		_ = g.NightOneVote(voteProvider, nil)
		break
	}

	return
}

func (g *Game) waitTwoVoteRoleFakeTimer(allPlayersWithRole []*playerPack.Player) {
	g.randomTimer()

	select {
	case voteP := <-g.TwoVoteChan:
		// All votes will be with errors
		err := g.NightTwoVote(voteP, nil)
		g.infoSender <- newErrSignal(err)
	case <-g.timerDone:
		votedPlayer := int(allPlayersWithRole[0].ID)
		voteProvider := NewTwoVoteProvider(strconv.Itoa(votedPlayer), EmptyVoteStr, EmptyVoteStr, false)
		_ = g.NightTwoVote(voteProvider, nil)
		break
	}

	for _, p := range allPlayersWithRole {
		p.Votes = append(p.Votes, EmptyVoteInt)
	}
}

func (g *Game) twoVoterRoleNightVoting(allPlayersWithRole []*playerPack.Player,
	containsNotMutedPlayers bool, deadline time.Duration) {
	var err error

	if !containsNotMutedPlayers {
		g.waitTwoVoteRoleFakeTimer(allPlayersWithRole)
		return
	}

	g.timer(deadline)

	select {
	case voteP := <-g.TwoVoteChan:
		err = g.NightTwoVote(voteP, nil)
		if err == nil {
			g.timerStop <- struct{}{}
			break
		} else {
			g.infoSender <- newErrSignal(err)
		}
	case <-g.timerDone:
		votedPlayer := int(allPlayersWithRole[0].ID)
		voteProvider := NewTwoVoteProvider(strconv.Itoa(votedPlayer), EmptyVoteStr, EmptyVoteStr, false)
		_ = g.NightTwoVote(voteProvider, nil)
		break
	}

	return
}

// AffectNight changes players according to the night's actions.
// Errors during execution are sent to the channel
func (g *Game) AffectNight(l NightLog) {
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
						}
					}
					select {
					case <-g.ctx.Done():
						return
					default:
						safeSendErrSignal(g.infoSender, channelPack.FromUserToSpectator(mainChannel, p.Tag))
					}
				}
			}
		}(newDeadPersons)

		// Changing arrays according to the night
		g.Active = &newActivePlayers
		g.Dead.Add(newDeadPersons...)
		g.Unlock()

		// Sending a message about who died today.
		err := g.Messenger.AfterNight.SendAfterNightMessage(l, g.MainChannel)
		safeSendErrSignal(g.infoSender, err)
		// Then, for each person try to do his reincarnation
		for _, p := range newActivePlayers {
			g.reincarnation(p)
		}
		return
	}
}
