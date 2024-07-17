package game

import (
	"math"
	"strconv"
	"time"

	"github.com/https-whoyan/MafiaBot/core/player"
)

const (
	DayPersentageToNextStage = 50
)

func (g *Game) Day() DayLog {
	select {
	case <-g.ctx.Done():
		return DayLog{}
	default:
		g.SetState(DayState)
		g.infoSender <- g.newSwitchStateSignal()

		g.RLock()
		deadline := CalculateDayDeadline(
			g.NightCounter, g.Dead.Len(), g.RolesConfig.PlayersCount)
		g.RUnlock()
		safeSendErrSignal(g.infoSender, g.Messenger.Day.SendMessageAboutNewDay(g.MainChannel, deadline))

		return g.StartDayVoting(deadline)
	}
}

func (g *Game) StartDayVoting(deadline time.Duration) DayLog {
	votesMp := make(map[int]int)
	occurrencesMp := make(map[int]int)

	g.timer(deadline)

	var kickedID = -1
	var breakDownDayPlayersCount = int(math.Ceil(float64(DayPersentageToNextStage*g.Active.Len()) / 100.0))

	acceptTheVote := func(voteP VoteProviderInterface) (kickedID *int) {
		var votedPlayerID = int(g.Active.SearchPlayerByID(voteP.GetVotedPlayerID()).ID)
		var vote, _ = strconv.Atoi(voteP.GetVote())

		if prevVote, isContains := votesMp[votedPlayerID]; isContains {
			occurrencesMp[prevVote]--
		}
		occurrencesMp[vote]++
		votesMp[votedPlayerID] = vote

		// If occurrencesMp[vote] >= breakDownDayPlayersCount
		if occurrencesMp[vote] >= breakDownDayPlayersCount {
			kickedID = &vote
		}
		// Case, when all players leave his vote
		if len(votesMp) == g.Active.Len() {
			// Calculate pVote, which have maximum occurrences
			var (
				mxOccurrence = 0
				mxVote       = 0
			)

			for pVote, occurrence := range occurrencesMp {
				if occurrence > mxOccurrence {
					mxOccurrence = occurrence
					mxVote = pVote
				}
			}

			return &mxVote
		}

		return
	}

	dayLog := DayLog{
		DayNumber: g.NightCounter,
		IsSkip:    true,
	}

	standDayLog := func(kickedID *int) {
		dayLog.Kicked = kickedID
		dayLog.DayVotes = votesMp
		dayLog.IsSkip = false
		if *kickedID == -1 {
			dayLog.Kicked = nil
			dayLog.IsSkip = true
		}
	}

	for {
		isNeedToContinue := false
		select {
		case <-g.ctx.Done():
			break
		case <-g.timerDone:
			break
		case voteP := <-g.VoteChan:
			err := g.DayVote(voteP, nil)
			if err != nil {
				isNeedToContinue = true
				break
			}
			maybeKickedID := acceptTheVote(voteP)
			if maybeKickedID != nil {
				kickedID = *maybeKickedID
				g.timerDone <- struct{}{}
				break
			}
		}
		if !isNeedToContinue {
			break
		}
	}

	standDayLog(&kickedID)
	return dayLog
}

// CalculateDayDeadline calculate the day max time.
func CalculateDayDeadline(nighCounter int, deadCount int, totalPlayers int) time.Duration {
	// Weights of aspects
	const (
		currNightCounterWeight  = 0.61
		deadCountWeight         = 0.68
		totalPlayersCountWeight = 0.27
	)

	var basicMinutes = 0.0 // TODO 2.2
	nightCounterAddMinutes := currNightCounterWeight * float64(nighCounter)
	deadCountAddMinutes := deadCountWeight * float64(deadCount)
	totalPlayersCountAddMinutes := totalPlayersCountWeight * float64(totalPlayers)

	totalTime := basicMinutes + nightCounterAddMinutes + deadCountAddMinutes + totalPlayersCountAddMinutes
	totalTimeMinutes := math.Ceil(totalTime)
	return time.Minute * time.Duration(totalTimeMinutes)
}

func (g *Game) AffectDay(l DayLog) (isFool bool) {
	if l.IsSkip {
		safeSendErrSignal(g.infoSender, g.Messenger.Day.SendMessageThatDayIsSkipped(g.MainChannel))
		return
	}
	kickedPlayer := (*g.Active)[player.IDType(*l.Kicked)]
	safeSendErrSignal(g.infoSender, g.Messenger.Day.SendMessageAboutKickedPlayer(g.MainChannel, kickedPlayer))

	g.Active.ToDead(kickedPlayer.ID, player.KilledByDayVoting, g.NightCounter, g.Dead)
	return
}
