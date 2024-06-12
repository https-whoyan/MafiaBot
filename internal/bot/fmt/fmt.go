package fmt

// _____________
// Text Style, FMTer
// _____________

// DiscordFMTer FmtInterface core realization

var FMTInstance = &DiscordFMTer{}

type DiscordFMTer struct{}

func (DiscordFMTer) Bold(s string) string {
	return "**" + s + "**"
}
func (DiscordFMTer) Italic(s string) string {
	return "_" + s + "_"
}
func (DiscordFMTer) Underline(s string) string {
	return "__" + string(s) + "__"
}

// CodeBlock language may be empty, it's ok
func (DiscordFMTer) CodeBlock(language string, s string) string {
	return "```" + language + s + "```"
}

// For less code.

func (f DiscordFMTer) B(s string) string {
	return f.Bold(s)
}
func (f DiscordFMTer) I(s string) string {
	return f.Italic(s)
}
func (f DiscordFMTer) U(s string) string {
	return f.Underline(s)
}
func (f DiscordFMTer) CD(s string) string {
	return f.CodeBlock("", s)
}

// B + U / B + I / I + U

func (f DiscordFMTer) BU(s string) string {
	return f.B(f.U(s))
}
func (f DiscordFMTer) BI(s string) string {
	return f.B(f.I(s))
}
func (f DiscordFMTer) IU(s string) string {
	return f.I(f.U(s))
}

// __________
// Stickers
// __________

var (
	RegistrationPlayerSticker    = ":grin:"
	RegistrationSpectatorSticker = ":smiling_imp:"
	LuckySticker                 = ":four_leaf_clover:"
)

// MappedStickersUnicode save stickers Unicode:
var MappedStickersUnicode = map[string]string{
	RegistrationPlayerSticker:    "U+1F601",
	RegistrationSpectatorSticker: "U+1F608",
}

func GetUnicodeBySticker(sticker string) string {
	return MappedStickersUnicode[sticker]
}
