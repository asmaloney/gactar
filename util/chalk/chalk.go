// Package chalk provides convenience utilities for working with colourized output using gchalk.
package chalk

import (
	"fmt"

	"github.com/jwalton/gchalk"
)

var (
	Error   = gchalk.WithBold().Red
	Warning = gchalk.Yellow

	Bold   = gchalk.Bold
	Header = gchalk.Green
)

func PrintErr(err error) {
	fmt.Println(Error("error:", err.Error()))
}

func PrintErrStr(str ...string) {
	fmt.Print(Error("error: "))
	fmt.Println(Error(str...))
}

func PrintWarningStr(str ...string) {
	fmt.Println(Warning(str...))
}
