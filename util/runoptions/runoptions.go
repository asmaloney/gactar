// Package runoptions preovide structs and types to pass around options used when running models.
package runoptions

import (
	"slices"

	"github.com/asmaloney/gactar/util/container"
)

type ACTRLogLevel string

var (
	// ValidFrameworks lists the valid options for choosing frameworks on the command line and in the
	// interactive case. Make sure "all" is the first entry as we use [1:] to get the rest.
	ValidFrameworks = []string{"all", "ccm", "pyactr", "vanilla"}

	ACTRLoggingLevels = []string{
		"min",
		"info",
		"detail",
	}
)

// ValidLogLevel returns whether the string is a valid logging level or not.
func ValidLogLevel(e string) bool {
	return slices.Contains(ACTRLoggingLevels, e)
}

// FrameworkNameList is a list of framework names used in the run options.
type FrameworkNameList []string

// InitialBuffers is a map of buffer names to initial contents of the buffer.
// This is used when passing in user-defined initial contents e.g. through a web API.
type InitialBuffers map[string]string

// Options are options used when running a model.
type Options struct {
	// List of frameworks to run on (if empty, "all" which means all active frameworks)
	// This is used in web mode to let the user select which frameworks to run on.
	// With the CLI options, this will be always be "all" since they specify the
	// frameworks on the command line.
	Frameworks FrameworkNameList

	// Stores the initial contents of any buffers
	InitialBuffers InitialBuffers

	// One of 'min', 'info', or 'detail'
	LogLevel *ACTRLogLevel

	// If true, output detailed info about activations
	TraceActivations *bool

	// The seed to use for generating pseudo-random numbers (allows for reproducible runs)
	// For all frameworks, if it is not set it uses current system time.
	// Use a uint32 because pyactr uses numpy and that's what its random number seed uses.
	RandomSeed *uint32
}

// New returns a default-initialized Options struct.
func New() Options {
	logLevel := ACTRLogLevel("info")
	activations := false

	return Options{
		Frameworks:       FrameworkNameList{"all"},
		LogLevel:         &logLevel,
		TraceActivations: &activations,
		RandomSeed:       nil,
	}
}

// IsValidFramework returns if the framework name is in our list of valid ones or not.
func IsValidFramework(frameworkName string) bool {
	return slices.Contains(ValidFrameworks, frameworkName)
}

// ValidNamedFrameworks returns the list of all valid framework names without "all".
func ValidNamedFrameworks() []string {
	return ValidFrameworks[1:]
}

// NormalizeFrameworkList will look for "all" and replace it with all available
// framework names.
func (f *FrameworkNameList) NormalizeFrameworkList(activeFrameworks FrameworkNameList) {
	if f == nil || slices.Contains(*f, "all") {
		*f = activeFrameworks
	}

	*f = container.UniqueAndSorted(*f)
}

// VerifyFrameworkList will check that each name is of a valid framework and that
// it is active on this server.
func (f FrameworkNameList) VerifyFrameworkList(activeFrameworks FrameworkNameList) (err error) {
	for _, name := range f {
		if !IsValidFramework(name) {
			err = &ErrInvalidFrameworkName{Name: name}
			return
		}

		// we have a valid name, check if it is active
		if !slices.Contains(activeFrameworks, name) {
			err = &ErrFrameworkNotActive{Name: name}
			return
		}
	}

	return
}
