package fmt

// _____________
// Text Style, FMTer
// _____________

var DiscordFMTInstance = &DiscordFMTer{}

// DiscordFMTer FmtInterface core realization
type DiscordFMTer struct{}

func (DiscordFMTer) Bold(s string) string      { return "**" + s + "**" }
func (DiscordFMTer) Italic(s string) string    { return "_" + s + "_" }
func (DiscordFMTer) Underline(s string) string { return "__" + string(s) + "__" }
func (DiscordFMTer) Block(s string) string     { return "`" + s + "`" }
func (DiscordFMTer) LineSplitter() string      { return "\n" }
func (DiscordFMTer) InfoSplitter() string      { return "==============================" }
func (DiscordFMTer) Tab() string               { return "\t" }
func (DiscordFMTer) Mention() string           { return "@" }

// For less code.

func (f DiscordFMTer) B(s string) string  { return f.Bold(s) }
func (f DiscordFMTer) I(s string) string  { return f.Italic(s) }
func (f DiscordFMTer) U(s string) string  { return f.Underline(s) }
func (f DiscordFMTer) Bl(s string) string { return f.Block(s) }
func (f DiscordFMTer) NL() string         { return f.LineSplitter() }
func (f DiscordFMTer) M() string          { return f.Mention() }

// B + U / B + I / I + U

func (f DiscordFMTer) BU(s string) string { return f.B(f.U(s)) }
func (f DiscordFMTer) BI(s string) string { return f.B(f.I(s)) }
func (f DiscordFMTer) IU(s string) string { return f.I(f.U(s)) }

// __________
// Stickers
// __________

var (
	RegistrationPlayerSticker    = "😁"
	RegistrationSpectatorSticker = "😈"
	ConfigChoiceSticker          = "👍🏻"
)
