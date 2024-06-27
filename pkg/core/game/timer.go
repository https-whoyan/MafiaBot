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

// ParalleledVoteTimer Init voteTimer in own goroutine
func ParalleledVoteTimer(ch chan<- VoteProviderInterface, done <-chan struct{},
	duration time.Duration, votingUserID string, isServerUserID bool, wg *sync.WaitGroup) {
	go voteTimer(ch, done, duration, votingUserID, isServerUserID, wg)
}

// VoteFakeTimer is used to simulate the operation of a timer.
// Selects a random time from the Range and sends an empty voice after it has passed
func voteFakeTimer(ch chan<- VoteProviderInterface, votingUserID string, isServerUserID bool) {
	minMilliSecond := myTime.FakeVotingMinSeconds * 1000
	maxMilliSecond := myTime.FakeVotingMaxSeconds * 1000
	randMilliSecondDuration := rand.Intn(maxMilliSecond-minMilliSecond+1) + minMilliSecond
	duration := time.Duration(randMilliSecondDuration) * time.Millisecond

	voteProvider := NewVoteProvider(votingUserID, EmptyVoteStr, isServerUserID)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		ch <- voteProvider
	}
}

// ParalleledFakeTimer start voteFakeTimer in own goroutine
func ParalleledFakeTimer(ch chan<- VoteProviderInterface, votingUserID string, isServerUserID bool) {
	go voteFakeTimer(ch, votingUserID, isServerUserID)
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
	go twoVoteTimer(ch, done, duration, votingUserID, isServerUserID, wg)
}

// voteTwoFakeTimer is used to simulate the operation of a timer.
// Selects a random time from the Range and sends an empty voices after it has passed
func voteTwoFakeTimer(ch chan<- TwoVoteProviderInterface, votingUserID string, isServerUserID bool) {
	minMilliSecond := myTime.FakeVotingMinSeconds * 1000
	maxMilliSecond := myTime.FakeVotingMaxSeconds * 1000
	randMilliSecondDuration := rand.Intn(maxMilliSecond-minMilliSecond+1) + minMilliSecond
	duration := time.Duration(randMilliSecondDuration) * time.Millisecond

	voteProvider := NewTwoVoteProvider(votingUserID, EmptyVoteStr, EmptyVoteStr, isServerUserID)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		ch <- voteProvider
	}
}

// TwoVoteFakeTimer start voteTwoFakeTimer in own goroutine
func TwoVoteFakeTimer(ch chan<- TwoVoteProviderInterface, votingUserID string, isServerUserID bool) {
	go voteTwoFakeTimer(ch, votingUserID, isServerUserID)
}

// ____________________________________________
// Used to simulates that the role is alive.
// _____________________________________________

func FullFakeVoteTimer(ch chan<- VoteProviderInterface) {
	sleepRandomSecond()
	ch <- nil
}

func FullFakeTwoVotesTimer(ch chan<- TwoVoteProviderInterface) {
	sleepRandomSecond()
	ch <- nil
}

func sleepRandomSecond() {
	minMilliSecond := myTime.FakeVotingMinSeconds * 1000
	maxMilliSecond := myTime.FakeVotingMaxSeconds * 1000
	randMilliSecondDuration := rand.Intn(maxMilliSecond-minMilliSecond+1) + minMilliSecond
	duration := time.Duration(randMilliSecondDuration) * time.Millisecond

	time.Sleep(duration)
}
