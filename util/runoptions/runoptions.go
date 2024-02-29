// Package runoptions preovide structs and types to pass around options used when running models.
package runoptions

import "slices"

type ACTRLogLevel string

var ACTRLoggingLevels = []string{
	"min",
	"info",
	"detail",
}

// ValidLogLevel returns whether the string is a valid logging level or not.
func ValidLogLevel(e string) bool {
	return slices.Contains(ACTRLoggingLevels, e)
}

// Options are options used when running a model.
type Options struct {
	// One of 'min', 'info', or 'detail'
	LogLevel ACTRLogLevel

	// Output detailed info about activations
	TraceActivations bool

	// The seed to use for generating pseudo-random numbers (allows for reproducible runs)
	// For all frameworks, if it is not set it uses current system time.
	// Use a uint32 because pyactr uses numpy and that's what its random number seed uses.
	RandomSeed *uint32
}

// New returns a default-initialized Options struct.
func New() Options {
	return Options{
		LogLevel:         ACTRLogLevel("info"),
		TraceActivations: false,
		RandomSeed:       nil,
	}
}
