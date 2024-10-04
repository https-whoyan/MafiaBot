package time

// Time and deadline constants are described below
const (
	RegistrationDeadlineSeconds     = 300
	VotingGameConfigDeadlineSeconds = 300
)

// Everything below is automatically calculated
const (
	RegistrationDeadlineMinutes     = RegistrationDeadlineSeconds / 60
	VotingGameConfigDeadlineMinutes = VotingGameConfigDeadlineSeconds / 60
)

const (
	BotTimeFormat = "2006-01-02 15:04:05"
)
