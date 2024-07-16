package game

import (
	"context"
	"log"
	"time"

	"github.com/https-whoyan/MafiaBot/core/roles"
)

// Signal is the structure that is used as the send data, when the game starts.
// Reports game state changes as well as errors that occur during game play
//
// Example: HandleExample
type Signal interface {
	GetSignalType() SignalType
	GetTime() time.Time
}

type SignalType int8

const (
	SwitchStateSignalType SignalType = 0
	ErrorSignalType       SignalType = 1
	CloseSignalSignalType SignalType = 2
)

// _____________________
// SwitchStateSignal
// _____________________

// SwitchStateSignal - the signal returned when the stage of the game has changed.
//
// See the declaration
type SwitchStateSignal struct {
	InitialTime      time.Time
	SignalType       SignalType
	SwitchSignalType SwitchStateType
	Value            SwitchValue
}

func (s SwitchStateSignal) GetSignalType() SignalType { return s.SignalType }
func (s SwitchStateSignal) GetTime() time.Time        { return s.InitialTime }

type SwitchStateType int8

const (
	SwitchStateSwitchStateType          SwitchStateType = 0
	SwitchNightVotedRoleSwitchStateType SwitchStateType = 1
)

// SwitchValue is used for polymorphism.
// Depending on the SwitchStateType, it will return either SwitchStateValue or SwitchNightVoteRoleSwitchValue
type SwitchValue interface {
	GetValue() SwitchValue
}

type SwitchStateValue struct {
	PreviousValue State
	CurrentState  State
}

func (s SwitchStateValue) GetValue() SwitchValue { return s }

type SwitchNightVoteRoleSwitchValue struct {
	CurrentVotedRole *roles.Role
	IsTwoVotes       bool
}

func (s SwitchNightVoteRoleSwitchValue) GetValue() SwitchValue { return s }

// _____________________
// ErrSignal
// _____________________

type ErrSignal struct {
	InitialTime   time.Time
	SignalType    SignalType
	ErrSignalType ErrSignalType
	Err           error
}

func (s ErrSignal) GetSignalType() SignalType { return s.SignalType }
func (s ErrSignal) GetTime() time.Time        { return s.InitialTime }

type ErrSignalType int

const (
	ErrorSignal ErrSignalType = 0
	FatalSignal ErrSignalType = 1 // After a fatal signal, the channel will close immediately.
)

// _____________________
// Close signal
// _____________________

type CloseSignal struct { // After a close signal, the channel will close immediately.
	InitialTime time.Time
	SignalType  SignalType
	Message     string
}

func (c CloseSignal) GetSignalType() SignalType { return c.SignalType }
func (c CloseSignal) GetTime() time.Time        { return c.InitialTime }

// _________________
// Internal code.
//__________________

func (g *Game) newSwitchStateSignal() Signal {
	return SwitchStateSignal{
		SignalType:       SwitchStateSignalType,
		InitialTime:      time.Now(),
		SwitchSignalType: SwitchStateSwitchStateType,
		Value: SwitchStateValue{
			CurrentState:  g.State,
			PreviousValue: g.PreviousState,
		},
	}
}

func (g *Game) newSwitchVoteSignal() Signal {
	return SwitchStateSignal{
		SignalType:       SwitchStateSignalType,
		InitialTime:      time.Now(),
		SwitchSignalType: SwitchNightVotedRoleSwitchStateType,
		Value: SwitchNightVoteRoleSwitchValue{
			CurrentVotedRole: g.NightVoting,
			IsTwoVotes:       g.NightVoting.IsTwoVotes,
		},
	}
}

func newErrSignal(err error) Signal {
	return ErrSignal{
		SignalType:    ErrorSignalType,
		InitialTime:   time.Now(),
		ErrSignalType: ErrorSignal,
		Err:           err,
	}
}

func newFatalSignal(err error) Signal {
	return ErrSignal{
		SignalType:    ErrorSignalType,
		InitialTime:   time.Now(),
		ErrSignalType: FatalSignal,
		Err:           err,
	}
}

func newCloseSignal(msg string) Signal {
	return CloseSignal{
		InitialTime: time.Now(),
		SignalType:  CloseSignalSignalType,
		Message:     msg,
	}
}

func safeSendErrSignal(ch chan<- Signal, err error) {
	if err == nil {
		return
	}
	ch <- newErrSignal(err)
}

func sendFatalSignal(ch chan<- Signal, err error) {
	fatalSignal := newFatalSignal(err)
	ch <- fatalSignal
	close(ch)
}

func sendCloseSignal(ch chan<- Signal, msg string) {
	ch <- newCloseSignal(msg)
	close(ch)
}

// ____________________
// Example
// ____________________

// HandleExample Example of a handler
func HandleExample() {
	var g = &Game{}
	var ctx context.Context // Your context (prefer Background)
	// Suppose you have initialized the game. (Called the Init method)
	var signalChannel <-chan Signal
	signalChannel = g.Run(ctx)

	for signal := range signalChannel {

		if signal.GetSignalType() == SwitchStateSignalType {
			switchSignal := signal.(SwitchStateSignal)

			if switchSignal.SwitchSignalType == SwitchStateSwitchStateType {

				currentGameState := switchSignal.Value.(SwitchStateValue).CurrentState
				log.Println(currentGameState) // For no errors
				// currentGameState now indicates which stage the game has switched to.

			} else if switchSignal.SwitchSignalType == SwitchNightVotedRoleSwitchStateType {

				switchVoteRoleValue := switchSignal.Value.(SwitchNightVoteRoleSwitchValue)
				currVotedRole, isTwoVotes := switchVoteRoleValue.CurrentVotedRole, switchVoteRoleValue.IsTwoVotes
				log.Println(currVotedRole, isTwoVotes) // For no errors
				// currVotedRole now indicates which role should Vote now.
				// isTwoVotes indicates whether 2 voices for a role are used at once or not.
				// This will help you understand which channel you need to send voice data to.
				//
				// Also, this information can can be learned from the structure of the game - NighVoting.

			}

		} else if signal.GetSignalType() == ErrorSignalType {
			errSignal := signal.(ErrSignal)

			errSignalType := errSignal.ErrSignalType
			err := errSignal.Err // Error
			if errSignalType == FatalSignal {
				// This means that the game will immediately end and the channel will be closed.
				// The FinishAnyway method will be called automatically.
				log.Println(err)
			} else if errSignalType == ErrorSignal {
				// You got error, please, process it.
				log.Println(err)
			}

		} else if signal.GetSignalType() == CloseSignalSignalType {
			// This type means that the game has been successfully played, a message has been sent to the channel,
			// and the channel will close after this message.
			message := signal.(CloseSignal).Message
			log.Println(message)
		}
	}
}
