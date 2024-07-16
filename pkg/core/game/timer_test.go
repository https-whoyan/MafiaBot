package game

// TODO
// REPLACE

/*

var (
	senderUserFakeID    = "1"
	senderIDIsServerID  = false
	receiverUserFakeID  = "2"
	receiverUserFakeID2 = "3"
)

func TestTimer(t *testing.T) {
	t.Parallel()
	t.Run("Test1", func(t *testing.T) {
		t.Parallel()
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})
		duration := 5 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			ch <- resaverVoteProvider
			close(done)
		}()
		vote := (<-ch).GetVote()
		wg.Wait()
		if vote != receiverUserFakeID {
			t.Errorf("Got Vote %v, expected receiverUserFakeID %v", vote, receiverUserFakeID)
		}
	})

	t.Run("Test2", func(t *testing.T) {
		t.Parallel()
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})
		duration := 5 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			time.Sleep(4 * time.Second)
			ch <- resaverVoteProvider
			close(done)
		}()
		vote := (<-ch).GetVote()
		wg.Wait()
		if vote != receiverUserFakeID {
			t.Errorf("Got Vote %v, expected receiverUserFakeID %v", vote, receiverUserFakeID)
		}
	})

	t.Run("Test3", func(t *testing.T) {
		t.Parallel()
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})
		duration := 5 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			time.Sleep(6 * time.Second)
			ch <- resaverVoteProvider
			close(done)
		}()
		vote := (<-ch).GetVote()
		wg.Wait()
		if vote != EmptyVoteStr {
			t.Errorf("Got Vote %v, expected EmptyVoteStr", vote)
		}
	})

	t.Run("Test4", func(t *testing.T) {
		t.Parallel()
		ch := make(chan VoteProviderInterface)
		done := make(chan struct{})
		duration := 2 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewVoteProvider(senderUserFakeID, receiverUserFakeID, senderIDIsServerID)
		VoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			time.Sleep(4 * time.Second)
			ch <- resaverVoteProvider
			close(done)
		}()
		vote := (<-ch).GetVote()
		wg.Wait()
		if vote != EmptyVoteStr {
			t.Errorf("Got Vote %v, expected EmptyVoteStr", vote)
		}
	})
}

func TestTwoVoteTimer(t *testing.T) {
	t.Parallel()
	t.Run("Test1", func(t *testing.T) {
		t.Parallel()
		ch := make(chan TwoVoteProviderInterface)
		done := make(chan struct{})
		duration := 5 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewTwoVoteProvider(senderUserFakeID, receiverUserFakeID, receiverUserFakeID2, senderIDIsServerID)
		TwoVoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			ch <- resaverVoteProvider
			close(done)
		}()
		vote1, vote2 := (<-ch).GetVote()
		wg.Wait()
		if vote1 != receiverUserFakeID || vote2 != receiverUserFakeID2 {
			t.Errorf("Got Votes %v and %v, expected receiverUserFakeID %v and %v", vote1,
				vote2, receiverUserFakeID, receiverUserFakeID2)
		}
	})

	t.Run("Test2", func(t *testing.T) {
		t.Parallel()
		ch := make(chan TwoVoteProviderInterface)
		done := make(chan struct{})
		duration := 5 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewTwoVoteProvider(senderUserFakeID, receiverUserFakeID, receiverUserFakeID2, senderIDIsServerID)
		TwoVoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			time.Sleep(4 * time.Second)
			ch <- resaverVoteProvider
			close(done)
		}()
		vote1, vote2 := (<-ch).GetVote()
		wg.Wait()
		if vote1 != receiverUserFakeID || vote2 != receiverUserFakeID2 {
			t.Errorf("Got Votes %v and %v, expected receiverUserFakeID %v and %v", vote1,
				vote2, receiverUserFakeID, receiverUserFakeID2)
		}
	})

	t.Run("Test3", func(t *testing.T) {
		t.Parallel()
		ch := make(chan TwoVoteProviderInterface)
		done := make(chan struct{})
		duration := 5 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewTwoVoteProvider(senderUserFakeID, receiverUserFakeID, receiverUserFakeID2, senderIDIsServerID)
		TwoVoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			time.Sleep(6 * time.Second)
			ch <- resaverVoteProvider
			close(done)
		}()
		vote1, vote2 := (<-ch).GetVote()
		wg.Wait()
		if vote1 != EmptyVoteStr || vote2 != EmptyVoteStr {
			t.Errorf("Got Vote %v and %v, expected EmptyVoteStr's", vote1, vote2)
		}
	})

	t.Run("Test4", func(t *testing.T) {
		t.Parallel()
		ch := make(chan TwoVoteProviderInterface)
		done := make(chan struct{})
		duration := 2 * time.Second

		wg := &sync.WaitGroup{}
		wg.Add(2)

		resaverVoteProvider := NewTwoVoteProvider(senderUserFakeID, receiverUserFakeID, receiverUserFakeID2, senderIDIsServerID)
		TwoVoteTimer(ch, done, duration, senderUserFakeID, senderIDIsServerID, wg)
		go func() {
			defer wg.Done()
			time.Sleep(4 * time.Second)
			ch <- resaverVoteProvider
			close(done)
		}()
		vote1, vote2 := (<-ch).GetVote()
		wg.Wait()
		if vote1 != EmptyVoteStr || vote2 != EmptyVoteStr {
			t.Errorf("Got Vote %v %v, expected EmptyVoteStr", vote1, vote2)
		}
	})
}

func TestFakeTimer(t *testing.T) {
	t.Parallel()
	t.Run("Test1", func(t *testing.T) {
		t.Parallel()
		startTime := time.Now()
		done := make(chan struct{})
		FakeTimer(done)
		<-done
		close(done)
		endTime := time.Now()
		log.Print(
			"Fake timer runs: ",
			float64(endTime.Sub(startTime).Milliseconds())/1000.0,
			" seconds")
	})
}

*/
