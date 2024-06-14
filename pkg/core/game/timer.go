package game

import (
	"math/rand"
	"time"

	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// _______________________
// VoteTimers
// _______________________

// VoteTimer sends an empty voice to the transmitted channel if duration have elapsed.
func VoteTimer(ch chan<- VoteProviderInterface, done <-chan struct{},
	duration time.Duration, votingUserID string, isServerUserID bool) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

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

// VoteFakeTimer is used to simulate the operation of a timer.
// Selects a random time from the Range and sends an empty voice after it has passed
func VoteFakeTimer(ch chan<- VoteProviderInterface, votingUserID string, isServerUserID bool) {
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
