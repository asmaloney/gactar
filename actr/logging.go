package actr

type ACTRLogLevel string

var ACTRLoggingLevels = []string{
	"min",
	"info",
	"detail",
}

func ValidLogLevel(e string) bool {
	for _, a := range ACTRLoggingLevels {
		if a == e {
			return true
		}
	}

	return false
}
