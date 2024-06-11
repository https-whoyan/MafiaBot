package fmt

// _____________
// Text Style
// _____________

// BotFMTer FmtInterface core realization
type BotFMTer struct{}

func NewBotFMTer() *BotFMTer {
	return &BotFMTer{}
}

func (BotFMTer) Bold(s string) string {
	return "**" + s + "**"
}

func (BotFMTer) Italic(s string) string {
	return "_" + s + "_"
}

func (BotFMTer) Underline(s string) string {
	return "__" + string(s) + "__"
}

// CodeBlock language may be empty, it's ok
func (BotFMTer) CodeBlock(language string, s string) string {
	return "```" + language + s + "```"
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
