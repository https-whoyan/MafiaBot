package fmt

type FmtInterface interface {
	Bold(s string) string
	Italic(b string) string
	Underline(b string) string
	CodeBlock(language string, b string) string
}

func BoldItalic(f FmtInterface, s string) string {
	return f.Italic(f.Bold(s))
}

func BoldUnderline(f FmtInterface, s string) string {
	return f.Underline(f.Bold(s))
}

func ItalicUnderline(f FmtInterface, s string) string {
	return f.Italic(f.Underline(s))
}
