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

func (g *Game) Day(ch chan<- Signal) DayLog {
	select {
	case <-g.ctx.Done():
		return DayLog{}
	default:
		g.Lock()
		g.SetState(DayState)
		g.Unlock()
		ch <- g.newSwitchStateSignal()

		g.RLock()
		deadline := CalculateDayDeadline(
			g.NightCounter, len(*g.Dead.ConvertToPlayers()), g.RolesConfig.PlayersCount)
		safeSendErrSignal(ch, g.messenger.Day.SendMessageAboutNewDay(g.MainChannel, deadline))
		g.RUnlock()

		// Start timer
		done := make(chan struct{})
		startTime(done, deadline)
		return g.StartDayVoting(done)
	}
}

func startTime(done chan<- struct{}, duration time.Duration) {
	go func() {
		time.Sleep(duration)
		close(done)
	}()
}

func (g *Game) StartDayVoting(done <-chan struct{}) DayLog {
	votesMp := make(map[int]int)
	occurrencesMp := make(map[int]int)

	var kickedID = -1
	var breakDownDayPlayersCount = int(math.Ceil(float64(DayPersentageToNextStage*g.Active.Len())/100.0)) + 1

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
			// Calculate vote, which have maximum occurrences
			var (
				mxOccurrence = 0
				mxVote       = 0
			)

			for vote, occurrence := range votesMp {
				if occurrence > mxOccurrence {
					mxOccurrence = occurrence
					mxVote = vote
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
	}

	select {
	case <-g.ctx.Done():
		standDayLog(&kickedID)
		return dayLog
	case <-done:
		standDayLog(&kickedID)
		return dayLog
	default:
		for voteP := range g.VoteChan {
			err := g.DayVote(voteP, nil)
			if err != nil {
				maybeKickedID := acceptTheVote(voteP)
				if maybeKickedID != nil {
					kickedID = *maybeKickedID
					break
				}
			}
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

	var basicMinutes = 2.2
	nightCounterAddMinutes := currNightCounterWeight * float64(nighCounter)
	deadCountAddMinutes := deadCountWeight * float64(deadCount)
	totalPlayersCountAddMinutes := totalPlayersCountWeight * float64(totalPlayers)

	totalTime := basicMinutes + nightCounterAddMinutes + deadCountAddMinutes + totalPlayersCountAddMinutes
	totalTimeMinutes := math.Ceil(totalTime)
	return time.Minute * time.Duration(totalTimeMinutes)
}

func (g *Game) AffectDay(l DayLog, ch chan<- Signal) (isFool bool) {
	if l.IsSkip {
		safeSendErrSignal(ch, g.messenger.Day.SendMessageThatDayIsSkipped(g.MainChannel))
		return
	}
	kickedPlayer := (*g.Active)[player.IDType(*l.Kicked)]
	safeSendErrSignal(ch, g.messenger.Day.SendMessageAboutKickedPlayer(g.MainChannel, kickedPlayer))

	safeSendErrSignal(ch, g.messenger.Day.SendMessageAboutKickedPlayer(g.MainChannel, kickedPlayer))
	g.Active.ToDead(kickedPlayer.ID, player.KilledByDayVoting, g.NightCounter, g.Dead)
	return
}
