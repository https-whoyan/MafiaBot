package time

import (
	"math/rand"
	"time"
)

func Timer(ch chan<- int, done <-chan struct{}, secondCount int) {
	duration := time.Duration(secondCount) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		select {
		case ch <- -1:
			return
		case <-done:
			return
		}
	case <-done:
		return
	}
}

func FakeTimer(ch chan<- int) {
	minMilliSecond := FakeVotingMinSeconds * 1000
	maxMilliSecond := FakeVotingMaxSeconds * 1000
	randMilliSecondDuration := rand.Intn(maxMilliSecond-minMilliSecond+1) + minMilliSecond
	duration := time.Duration(randMilliSecondDuration) * time.Millisecond

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		ch <- -1
	}
}
