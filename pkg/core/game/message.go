package game

import (
	"fmt"
	"github.com/https-whoyan/MafiaBot/core/roles"
	"io"
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

type Messenger struct {
	f          myFMT.FmtInterface
	Init       *initMessenger
	Night      *nightMessenger
	AfterNight *afterNightMessenger
	Finish     *finishMessenger
}

type primitiveMessenger struct {
	f myFMT.FmtInterface
	g *Game
}

func (m primitiveMessenger) sendMessage(msg string, writer io.Writer) error {
	_, err := writer.Write([]byte(msg))
	return err
}

func NewGameMessanger(f myFMT.FmtInterface, g *Game) *Messenger {
	base := &primitiveMessenger{f: f, g: g}
	return &Messenger{
		f:          f,
		Init:       &initMessenger{base},
		Night:      &nightMessenger{base},
		AfterNight: &afterNightMessenger{base},
		Finish:     &finishMessenger{base},
	}
}

// ____________
// Init
// ____________

type initMessenger struct {
	*primitiveMessenger
}

func (m initMessenger) SendStartMessage(writer io.Writer) error {
	var message string

	f := m.f

	nl := f.LineSplitter()
	dNl := nl + nl
	iL := f.InfoSplitter()

	message = f.Bold("Have a good day, " + getRandomPlayersCalling() + "!")
	message += dNl
	message += myFMT.BoldUnderline(f, "Today, our players:") + nl

	var aboutIDMessages []string
	activePlayers := *m.g.Active
	sort.Slice(activePlayers, func(i, j int) bool {
		return activePlayers[playerPack.IDType(j)].ID < activePlayers[playerPack.IDType(j)].ID
	})
	for _, player := range activePlayers {
		messageAboutPlayerID := f.Tab() + f.Bold(sCap(getRandomPlayerCalling())) + " " + f.Mention(player.ServerNick)
		messageAboutPlayerID += " with IDType in game " + f.Block(sInt(int(player.ID)))

		aboutIDMessages = append(aboutIDMessages, messageAboutPlayerID)
	}
	message += strings.Join(aboutIDMessages, nl)

	if len(*m.g.Spectators) != 0 {
		message += dNl
		message += "From behind the scenes to support us: "
		var spectatorMentions []string

		for _, spectator := range *m.g.Spectators {
			spectatorMentions = append(spectatorMentions, f.Mention(spectator.ServerNick))
		}
		message += strings.Join(spectatorMentions, ", ")
	}

	message += nl + iL + nl
	message += myFMT.ItalicUnderline(f, "Selected game configuration:") + nl
	message += m.g.RolesConfig.GetMessageAboutConfig(f)
	message += nl + iL + nl

	// Redo it if it false!!!!
	message += "A private message has been sent to each of you, you can find your ID and role in it."
	message += nl
	message += "Also, " + f.Bold("if you have an active night role, you have been added to special channels, where "+
		"you can send commands to the bot anonymously")
	if len(*m.g.Spectators) != 0 {
		message += f.Italic(" (but there's no hiding from observers))))")
	}
	message += "." + nl
	if m.g.renameMode != NotRenameModeMode {
		message += nl
		message += "Also, all participants have been prefixed with their IDs to make it more convenient for you."
	}
	message += dNl + iL + nl
	message += f.Bold("Welcome, welcome, welcome... Happy hunger games and the odds be ever in your favor! ") +
		f.Italic("(Or just have a good game!) 🍀")

	_, err := writer.Write([]byte(message))
	return err
}

// ____________
// Night
// ____________

type nightMessenger struct {
	*primitiveMessenger
}

func (m *nightMessenger) SendInitialNightMessage(w io.Writer) error {
	f := m.f
	message := f.Bold("Night №") + f.Block(strconv.Itoa(m.g.NightCounter)) + " is coming." + f.LineSplitter()
	message += fmt.Sprintf("On this night you are played by %v players.", len(*m.g.Active))
	return m.sendMessage(message, w)
}

func (m *nightMessenger) SendInvitingToVoteMessage(p *playerPack.Player, deadlineInSeconds int, w io.Writer) error {
	m.g.RLock()
	defer m.g.RUnlock()
	f := m.f
	message := f.Bold("Hello, " + f.Mention(p.ServerNick) + ". It's your turn to Vote.")
	message += f.LineSplitter()
	message += myFMT.BoldUnderline(f, fmt.Sprintf("Deadline: %v seconds.", deadlineInSeconds))
	return m.sendMessage(message, w)
}

func (m *nightMessenger) SendToPlayerThatIsMutedMessage(p *playerPack.Player, w io.Writer) error {
	m.g.RLock()
	defer m.g.RUnlock()

	message := "Oops.... someone was muted today!" + m.f.Mention(p.ServerNick) +
		" , just chill, bro."
	return m.sendMessage(message, w)
}

func (m *nightMessenger) SendThanksToMutedPlayerMessage(p *playerPack.Player, writer io.Writer) error {
	m.g.RLock()
	defer m.g.RUnlock()
	message := m.f.Bold(m.f.Mention(p.ServerNick) + ", always thanks!")
	return m.sendMessage(message, writer)
}

// ____________
// AfterNight
// ____________

type afterNightMessenger struct {
	*primitiveMessenger
}

// GetAfterNightMessage provide a message to main chat after game.
func (m afterNightMessenger) SendAfterNightMessage(l NightLog, w io.Writer) error {
	f := m.f
	message := myFMT.BoldItalic(f, "Dear citizens!") + f.LineSplitter()
	message += f.Bold("Today, we're losing")
	if len(l.Dead) == 0 {
		message += "....  " + myFMT.BoldUnderline(f, "Just our nerve cells...") + f.LineSplitter()
		message += f.Bold("Everyone survived.")
		return m.sendMessage(message, w)
	}
	message += " " + f.Block(strconv.Itoa(len(l.Dead))) +
		f.Bold(" people")
	var mentions []string
	idsSet := cnvPack.SliceToSet(l.Dead)
	for _, p := range *m.g.Active {
		if idsSet[int(p.ID)] {
			mentions = append(mentions, f.Mention(p.ServerNick))
		}
	}
	message += " which is to say: " + strings.Join(mentions, ", ")
	message += f.LineSplitter() + f.LineSplitter()
	message += f.Bold("Dear victims, you have " +
		strconv.Itoa(myTime.LastWordDeadlineMinutes) + " minute to say your angry.")
	return m.sendMessage(message, w)
}

// _____________________
// Team victory message
// _____________________

type finishMessenger struct {
	*primitiveMessenger
}

func (m finishMessenger) SendMessagesAboutEndOfGame(l FinishLog, w io.Writer) error {
	var message string
	f := m.f
	message = f.Bold("Our game has come to an end.") + f.LineSplitter()
	err := m.SendMessageAboutWinner(l, w)
	if err != nil {
		return err
	}
	err = m.sendMessage(message, w)
	return err
}

func (m finishMessenger) SendMessageAboutWinner(l FinishLog, w io.Writer) error {
	var message string

	if l.IsFool {
		message += m.getFoolWinnerMessage()
	} else {
		message += m.getTeamWinnerMessage(*l.WinnerTeam)
	}
	return m.sendMessage(message, w)
}

func (m finishMessenger) getTeamWinnerMessage(team roles.Team) string {
	//Todo
	return ""
}

func (m finishMessenger) getFoolWinnerMessage() string {
	//Todo
	return ""
}

// Todo
func (m finishMessenger) SendParticipantAboutMessage(w io.Writer) string {
	f := m.f
	var message string
	message = myFMT.BoldUnderline(f, "And the roles of the participants were:")

	allPartitions := m.g.Active
	allPartitions.Append(m.g.Dead.ConvertToPlayers())

	return message
}
