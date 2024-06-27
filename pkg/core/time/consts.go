package time

// Time and deadline constants are described below
const (
	VotingDeadline = 30

	FakeVotingMinSeconds = 5
	FakeVotingMaxSeconds = 25

	LastWordDeadline = 60
)

// Everything below is automatically calculated
const (
	LastWordDeadlineMinutes = LastWordDeadline / 60
)
