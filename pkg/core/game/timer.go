package game

import (
	"math/rand"
	"sync"
	"time"

	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// _______________________
// OneVoteTimers
// _______________________

// voteTimer sends an empty voice to the transmitted channel if duration have elapsed.
func voteTimer(ch chan<- VoteProviderInterface, done <-chan struct{},
	duration time.Duration, votingUserID string, isServerUserID bool, wg *sync.WaitGroup) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	defer wg.Done()

	voteProvider := NewVoteProvider(votingUserID, EmptyVoteStr, isServerUserID)

	select {
	case <-ticker.C:
		select {
		case ch <- voteProvider:
			return
		case <-done:
			return
		}
	case <-done:
		return
	}
}

// VoteTimer Init voteTimer in own goroutine
func VoteTimer(ch chan<- VoteProviderInterface, done <-chan struct{},
	duration time.Duration, votingUserID string, isServerUserID bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go voteTimer(ch, done, duration, votingUserID, isServerUserID, wg)
}

// _______________________
// TwoVotesTimer
// _______________________

// twoVoteTimer sends an empty votes to the transmitted channel if duration have elapsed.
func twoVoteTimer(ch chan<- TwoVoteProviderInterface, done <-chan struct{},
	duration time.Duration, votingUserID string, isServerUserID bool, wg *sync.WaitGroup) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	defer wg.Done()

	voteProvider := NewTwoVoteProvider(votingUserID, EmptyVoteStr, EmptyVoteStr, isServerUserID)

	select {
	case <-ticker.C:
		select {
		case ch <- voteProvider:
			return
		case <-done:
			return
		}
	case <-done:
		return
	}
}

// TwoVoteTimer Init twoVoteTimer in own goroutine
func TwoVoteTimer(ch chan<- TwoVoteProviderInterface, done <-chan struct{},
	duration time.Duration, votingUserID string, isServerUserID bool, wg *sync.WaitGroup) {
	wg.Add(1)
	go twoVoteTimer(ch, done, duration, votingUserID, isServerUserID, wg)
}

// ____________________________________________
// Used to simulate.
// _____________________________________________

func FakeTimer(done chan<- struct{}) { go fakeTimer(done) }

func fakeTimer(done chan<- struct{}) {
	sleepRandomSecond()
	done <- struct{}{}
	close(done)
}

func sleepRandomSecond() {
	minMilliSecond := myTime.FakeVotingMinSeconds * 1000
	maxMilliSecond := myTime.FakeVotingMaxSeconds * 1000
	randMilliSecondDuration := rand.Intn(maxMilliSecond-minMilliSecond+1) + minMilliSecond
	duration := time.Duration(randMilliSecondDuration) * time.Millisecond

	time.Sleep(duration)
}
