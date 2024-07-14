package handlers

import (
	"github.com/https-whoyan/MafiaBot/internal/fmt"
	"log"
	"time"

	"github.com/https-whoyan/MafiaBot/core/game"
)

// ProcessGameChan Used to process chan, and send addition info, about, what command need to send to the chat.
func ProcessGameChan(g *game.Game, f *fmt.DiscordFMTer, signalChannel <-chan game.Signal) {
	for signal := range signalChannel {
		log.Println(signal)
		switch signal.GetSignalType() {

		case game.ErrorSignalType:
			log.Println(signal.(game.ErrSignal).Err)

		case game.SwitchStateSignalType:

			switchSignal := signal.(game.SwitchStateSignal)
			log.Println("Switch State", switchSignal.Value)

			if switchSignal.SwitchSignalType == game.SwitchStateSwitchStateType {

				currentGameState := switchSignal.Value.(game.SwitchStateValue).CurrentState
				if currentGameState == game.DayState {
					// Timing for game send messages. Then, send message about command.\
					time.Sleep(time.Millisecond * 500)
					var message string
					message += "Use " + f.BU("/vote") + " command to leave a vote" + f.NL()
					message = "To vote for skipping, " + f.B("use -1")

					_, _ = g.MainChannel.Write([]byte(message))
				}
			} else {
				log.Println(switchSignal.Value.(game.SwitchNightVoteRoleSwitchValue).CurrentVotedRole)
			}

		}
	}
}
