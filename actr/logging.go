package actr

import "slices"

type ACTRLogLevel string

var ACTRLoggingLevels = []string{
	"min",
	"info",
	"detail",
}

func ValidLogLevel(e string) bool {
	return slices.Contains(ACTRLoggingLevels, e)
}
