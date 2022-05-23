package actr

import "github.com/asmaloney/gactar/util/container"

type ACTRLogLevel string

var ACTRLoggingLevels = []string{
	"min",
	"info",
	"detail",
}

func ValidLogLevel(e string) bool {
	return container.Contains(e, ACTRLoggingLevels)
}
