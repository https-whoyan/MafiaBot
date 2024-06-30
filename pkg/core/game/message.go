package game

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	cnvPack "github.com/https-whoyan/MafiaBot/core/converter"
	myFMT "github.com/https-whoyan/MafiaBot/core/fmt"
	playerPack "github.com/https-whoyan/MafiaBot/core/player"
	myTime "github.com/https-whoyan/MafiaBot/core/time"
)

// !!!!!!!!!!!!!!!!!!
// The use of all functions here assumes that you g.Init() has already been
// called and successfully executed without errors.
// !!!!!!!!!!!!!!!!!!

var (
	playersCalling = []string{"poopsies", "players", "ladies and gentlemen", "citizens"}
	playerCalling  = []string{"poops", "ancient", "modern", "member"}
)

var (
	sInt = func(s int) string { return strconv.Itoa(s) }
	sCap = func(s string) string {
		if len(s) <= 1 {
			return s
		}
		return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	}
	getRandomPlayersCalling = func() string { return playersCalling[rand.Intn(len(playersCalling))] }
	getRandomPlayerCalling  = func() string { return playerCalling[rand.Intn(len(playerCalling))] }
)

// ____________
// Init
// ____________

func (g *Game) GetStartMessage() string {
	var message string

	f := g.fmtEr

	nl := f.LineSplitter()
	dNl := nl + nl
	iL := f.InfoSplitter()

	message = f.Bold("Have a good day, " + getRandomPlayersCalling() + "!")
	message += dNl
	message += myFMT.BoldUnderline(f, "Today, our players:") + nl

	var aboutIDMessages []string
	activePlayers := g.Active
	sort.Slice(activePlayers, func(i, j int) bool {
		return activePlayers[i].ID < activePlayers[j].ID
	})
	for _, player := range activePlayers {
		messageAboutPlayerID := f.Tab() + f.Bold(sCap(getRandomPlayerCalling())) + " " + f.Mention(player.ServerNick)
		messageAboutPlayerID += " with ID in game " + f.Block(sInt(player.ID))

		aboutIDMessages = append(aboutIDMessages, messageAboutPlayerID)
	}
	message += strings.Join(aboutIDMessages, nl)

	if len(g.Spectators) != 0 {
		message += dNl
		message += "From behind the scenes to support us: "
		var spectatorMentions []string

		for _, spectator := range g.Spectators {
			spectatorMentions = append(spectatorMentions, f.Mention(spectator.ServerNick))
		}
		message += strings.Join(spectatorMentions, ", ")
	}

	message += nl + iL + nl
	message += myFMT.ItalicUnderline(f, "Selected game configuration:") + nl
	message += g.RolesConfig.GetMessageAboutConfig(f)
	message += nl + iL + nl

	// Redo it if it false!!!!
	message += "A private message has been sent to each of you, you can find your ID and role in it."
	message += nl
	message += "Also, " + f.Bold("if you have an active night role, you have been added to special channels, where "+
		"you can send commands to the bot anonymously")
	if len(g.Spectators) != 0 {
		message += f.Italic(" (but there's no hiding from observers))))")
	}
	message += "." + nl
	if g.renameMode != NotRenameModeMode {
		message += nl
		message += "Also, all participants have been prefixed with their IDs to make it more convenient for you."
	}
	message += dNl + iL + nl
	message += f.Bold("Welcome, welcome, welcome... Happy hunger games and the odds be ever in your favor! ") +
		f.Italic("(Or just have a good game!) 🍀")

	return message
}

// ____________
// Night
// ____________

func (g *Game) getInitialNightMessage() string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr
	message := f.Bold("Night №") + f.Block(strconv.Itoa(g.NightCounter)) + " is coming." + f.LineSplitter()
	message += fmt.Sprintf("On this night you are played by %v players.", len(g.Active))
	return message
}

func (g *Game) getInvitingMessageToVote(p *playerPack.Player, deadlineInSeconds int) string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr
	message := f.Bold("Hello, " + f.Mention(p.ServerNick) + ". It's your turn to Vote.")
	message += f.LineSplitter()
	message += myFMT.BoldUnderline(f, fmt.Sprintf("Deadline: %v seconds.", deadlineInSeconds))
	return message
}

func (g *Game) getMessageToPlayerThatIsMuted(p *playerPack.Player) string {
	g.RLock()
	defer g.RUnlock()
	f := g.fmtEr

	message := "Oops.... someone was muted today!" + f.Mention(p.ServerNick) +
		" , just chill, bro."
	return message
}

func (g *Game) getThanksToMutedPlayerMessage(p *playerPack.Player) string {
	g.RLock()
	defer g.RUnlock()
	message := g.fmtEr.Bold(g.fmtEr.Mention(p.ServerNick) + ", always thanks!")
	return message
}

// ____________
// AfterNight
// ____________

// GetAfterNightMessage provide a message to main chat after game.
func (g *Game) GetAfterNightMessage(l NightLog) string {
	f := g.fmtEr
	message := myFMT.BoldItalic(f, "Dear citizens!") + f.LineSplitter()
	message += f.Bold("Today, we're losing")
	if len(l.Dead) == 0 {
		message += "....  " + myFMT.BoldUnderline(f, "Just our nerve cells...") + f.LineSplitter()
		message += f.Bold("Everyone survived.")
		return message
	}
	message += " " + f.Block(strconv.Itoa(len(l.Dead))) +
		f.Bold(" people")
	var mentions []string
	idsSet := cnvPack.SliceToSet(l.Dead)
	for _, p := range g.StartPlayers {
		if idsSet[p.ID] {
			mentions = append(mentions, f.Mention(p.ServerNick))
		}
	}
	message += " which is to say: " + strings.Join(mentions, ", ")
	message += f.LineSplitter() + f.LineSplitter()
	message += f.Bold("Dear victims, you have " +
		strconv.Itoa(myTime.LastWordDeadlineMinutes) + " minute to say your angry.")
	return message
}

// _____________________
// Team victory message
// _____________________

func (g *Game) GetMessageAboutWinner(l FinishLog) string {
	if l.IsFool {
		return g.getFoolWinnerMessage()
	}

}

func (g *Game) getTeamWinnerMessage()

func (g *Game) getFoolWinnerMessage() string {

}

func (g *Game) getParticipantAboutMessage() string {
	f := g.fmtEr
	var message string
	message := myFMT.BoldUnderline(f, "And the roles of the participants were:")
}
