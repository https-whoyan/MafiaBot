package game

import (
	"math/rand"
	"time"

	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// Used to simulate.

func getRandomDuration() time.Duration {
	minMilliSecond := myTime.FakeVotingMinSeconds * 1000
	maxMilliSecond := myTime.FakeVotingMaxSeconds * 1000
	randMilliSecondDuration := rand.Intn(maxMilliSecond-minMilliSecond+1) + minMilliSecond
	return time.Duration(randMilliSecondDuration) * time.Millisecond
}

func (g *Game) timer(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			g.timerDone <- struct{}{}
			return
		case <-g.timerStop:
			return
		}
	}()
}

func (g *Game) randomTimer() {
	duration := getRandomDuration()
	g.timer(duration)
}
