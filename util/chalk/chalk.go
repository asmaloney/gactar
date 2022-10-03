// Package chalk provides convenience utilities for working with colourized output using gchalk.
package chalk

import (
	"fmt"
	"os"

	"github.com/jwalton/gchalk"
)

var (
	Error   = gchalk.WithBold().Red
	Warning = gchalk.Yellow

	Bold   = gchalk.Bold
	Header = gchalk.Green
)

func PrintErr(err error) {
	fmt.Fprintln(os.Stderr, Error("error:", err.Error()))
}

func PrintErrStr(str ...string) {
	fmt.Fprint(os.Stderr, Error("error: "))
	fmt.Fprintln(os.Stderr, Error(str...))
}

func PrintWarningStr(str ...string) {
	fmt.Println(Warning(str...))
}
