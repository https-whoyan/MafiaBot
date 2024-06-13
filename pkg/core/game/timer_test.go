package game

import (
	"log"
	"testing"
	"time"
)

var (
	senderUserFakeID   = "1"
	senderIDIsServerID = false
	receiverUserFakeID = "2"
)

func TestTimer(t *testing.T) {
	t.Parallel()
	t.Run("Test1", func(t *testing.T) {
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})

		duration := 5 * time.Second
		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		go VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID)
		go func() {
			ch <- resaverVoteProvider
			done <- struct{}{}
		}()

		vote := (<-ch).GetVote()
		if vote != receiverUserFakeID {
			t.Errorf("Got Vote %v, expected receiverUserFakeID %v", vote, receiverUserFakeID)
		}
	})

	t.Run("Test2", func(t *testing.T) {
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})

		duration := 5 * time.Second
		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		go VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID)
		go func() {
			time.Sleep(4 * time.Second)
			ch <- resaverVoteProvider
			done <- struct{}{}
		}()

		vote := (<-ch).GetVote()
		if vote != receiverUserFakeID {
			t.Errorf("Got Vote %v, expected receiverUserFakeID %v", vote, receiverUserFakeID)
		}
	})

	t.Run("Test3", func(t *testing.T) {
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})

		duration := 5 * time.Second
		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		go VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID)
		go func() {
			time.Sleep(6 * time.Second)
			ch <- resaverVoteProvider
			done <- struct{}{}
		}()

		vote := (<-ch).GetVote()
		if vote != EmptyVoteStr {
			t.Errorf("Got Vote %v, expected EmptyVoteStr", vote)
		}
	})

	t.Run("Test4", func(t *testing.T) {
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})

		duration := 2 * time.Second
		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		go VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID)
		go func() {
			time.Sleep(4 * time.Second)
			ch <- resaverVoteProvider
			done <- struct{}{}
		}()

		vote := (<-ch).GetVote()
		if vote != EmptyVoteStr {
			t.Errorf("Got Vote %v, expected EmptyVoteStr", vote)
		}
	})
}

func TestFakeTimer(t *testing.T) {
	t.Parallel()
	t.Run("Test1", func(t *testing.T) {
		ch := make(chan VoteProviderInterface)

		startTime := time.Now()
		go VoteFakeTimer(ch, senderUserFakeID, senderIDIsServerID)
		vote := (<-ch).GetVote()
		if vote != EmptyVoteStr {
			t.Errorf("Got Vote %v, expected EmptyVoteStr", vote)
		}
		endTime := time.Now()
		log.Print(
			"Fake timer runs: ",
			float64(endTime.Sub(startTime).Milliseconds())/1000.0,
			" seconds")
	})
}
