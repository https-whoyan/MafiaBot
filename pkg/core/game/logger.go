package game

import (
	"time"

	"github.com/https-whoyan/MafiaBot/core/channel"
	"github.com/https-whoyan/MafiaBot/core/player"
	"github.com/https-whoyan/MafiaBot/core/roles"
	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// GameLogger allows you to save game information.
//
// The implementation is thrown when the game is initialized Init,
// logs are automatically loaded and saved to the implementation for
// saving in the run and Finish methods.
type GameLogger interface {
	InitNewGame(g *Game) error
	SaveNightLog(g *Game, log NightLog) error
	SaveDayLog(g *Game, log DayLog) error
	SaveFinishLog(g *Game, log FinishLog) error
}

// ____________
// NightLog
// ____________

// NightLog saves all votes, as well as the IDs of those
// players who turned out to be dead based on the results of the night.
type NightLog struct {
	NightNumber int `json:"number"`
	// Key - ID of the voted player
	// Value - usually a vote, but in case the role uses 2 votes - 2 votes at once.
	NightVotes map[int][]int `json:"votes"`
	Dead       []int         `json:"dead"`
}

// NewNightLog Gives the log after nightfall.
// Panics if not called after night or during voting.
func (g *Game) NewNightLog() NightLog {
	if g.ctx == nil {
		panic("Game is not initialized")
	}
	select {
	case <-g.ctx.Done():
		return NightLog{}
	default:
		if g.State != NightState {
			panic("Inappropriate use not after overnight")
		}
		if g.NightVoting != nil {
			panic("the function is called during the night, not after it!")
		}

		g.RLock()
		defer g.RUnlock()

		nightNumber := g.NightCounter
		nightVotes := make(map[int][]int)
		for _, p := range g.Active {
			if p.Role.NightVoteOrder == -1 {
				continue
			}

			votes := []int{}
			n := len(p.Votes)
			if p.Role.IsTwoVotes {
				votes = []int{p.Votes[n-2], p.Votes[n-1]}
			} else {
				votes = []int{p.Votes[n-1]}
			}
			nightVotes[p.ID] = votes
		}
		var dead []int
		for _, p := range g.Active {
			if p.LifeStatus == player.Dead {
				dead = append(dead, p.ID)
			}
		}
		return NightLog{
			NightNumber: nightNumber,
			NightVotes:  nightVotes,
			Dead:        dead,
		}
	}
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
		defer g.Unlock()

		// Splitting arrays.
		var newActivePlayers []*player.Player
		var newDeadPersons []*player.Player

		for _, p := range g.Active {
			if p.LifeStatus == player.Dead {
				newDeadPersons = append(newDeadPersons, p)
			} else {
				newActivePlayers = append(newActivePlayers, p)
			}
		}

		// I will add add add all killed players after a minute of players after a minute of
		// players after a minute, so, using goroutine.
		go func(newDeadPersons []*player.Player) {
			duration := myTime.LastWordDeadline * time.Second
			time.Sleep(duration)
			if g.TryLock() {
				defer g.Unlock()
			}
			// I'm adding new dead players to the spectators in the channels (so they won't be so bored)
			for _, p := range newDeadPersons {
				for _, interactionChannel := range g.RoleChannels {
					safeSendErrSignal(ch, channel.FromUserToSpectator(interactionChannel, p.Tag))
				}
				safeSendErrSignal(ch, channel.FromUserToSpectator(g.MainChannel, p.Tag))
			}
		}(newDeadPersons)

		// Changing arrays according to the night
		g.Active = newActivePlayers
		g.Dead = append(g.Dead, newDeadPersons...)

		// Sending a message about who died today.
		message := g.GetAfterNightMessage(l)
		_, err := g.MainChannel.Write([]byte(message))
		safeSendErrSignal(ch, err)
		return
	}

}

// _______________
// DayLog
// _______________

type DayLog struct {
	DayNumber int `json:"number"`
	// Key - ID of the player who was voted for during the day to be excluded from the game
	// Value - number of votes.
	DayVotes map[int]int `json:"votes"`
	Kicked   int         `json:"kicked"`
	IsSkip   bool        `json:"isSkip"`
}

type FinishLog struct {
	WinnerTeam  roles.Team `json:"winnerTeam"`
	TotalNights int        `json:"totalNights"`
}
