// Package chalk provides convenience utilities for working with colourized output using gchalk.
package chalk

import (
	"fmt"
	"os"

	"github.com/jwalton/gchalk"
)

var (
	Default = gchalk.Reset

	Error     = gchalk.Red
	ErrorBold = gchalk.WithRed().Bold
	Warning   = gchalk.Yellow
	Success   = gchalk.WithGreen().Bold

	Bold   = gchalk.Bold
	Italic = gchalk.Italic
	Header = gchalk.Green

	BlueUnderline     = gchalk.WithBlue().Underline
	BlueBoldUnderline = gchalk.WithBlue().WithBold().Underline
)

func HasColor() bool {
	return gchalk.GetLevel() != gchalk.LevelNone
}

func PrintErrLight(err error) {
	fmt.Fprintln(os.Stderr, Error("error:", err.Error()))
}

func PrintErr(err error) {
	fmt.Fprintln(os.Stderr, ErrorBold("error:", err.Error()))
}

func PrintErrStr(str ...string) {
	fmt.Fprint(os.Stderr, ErrorBold("error: "))
	fmt.Fprintln(os.Stderr, ErrorBold(str...))
}

func PrintWarningStr(str ...string) {
	fmt.Println(Warning(str...))
}

func QuotedItalic(text string) string {
	return gchalk.Italic(fmt.Sprintf("%q", text))
}
