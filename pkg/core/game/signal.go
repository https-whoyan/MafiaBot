package game

import "errors"

// Signal is the structure that is used as the send data, when the game starts.
//
// It broadcasts only errors that occur during the game.
type Signal struct {
	Type  SignalType
	Value any
}

type SignalType int

const (
	ErrSignal   SignalType = 1
	FatalSignal SignalType = 2
	CloseSignal SignalType = 2
)

var (
	FatalGameErr = errors.New("fatal error")
)

type SignalValue int

func NewErrSignal(err error) Signal   { return Signal{Type: ErrSignal, Value: err} }
func NewFatalSignal(err error) Signal { return Signal{Type: FatalSignal, Value: err} }
func NewCloseSignal(err error) Signal { return Signal{Type: CloseSignal, Value: err} }

// Used to send ErrSignal only if err != nil.
// Used for smaller code.
func safeSendErrSignal(ch chan<- Signal, err error) {
	if err != nil {
		ch <- NewErrSignal(err)
	}
}
