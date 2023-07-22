package actr

import "golang.org/x/exp/slices"

type ACTRLogLevel string

var ACTRLoggingLevels = []string{
	"min",
	"info",
	"detail",
}

func ValidLogLevel(e string) bool {
	return slices.Contains(ACTRLoggingLevels, e)
}
