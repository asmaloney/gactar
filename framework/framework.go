package framework

import (
	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/util/container"
)

// ValidFrameworks lists the valid options for choosing frameworks on the command line and in the
// interactive case.
var ValidFrameworks = []string{"all", "ccm", "pyactr", "vanilla"}

// Info provides basic info to set up a framework.
type Info struct {
	Name     string // name of the framework
	Language string // language the framework uses

	FileExtension string // file extension of the intermediate file

	ExecutableName string // name of the executable to run

	PythonRequiredPackages []string // (Python only) List of packages this framework requires
}

// RunResult is the result of a Run() call which runs the code using the framework's executable.
type RunResult struct {
	FileName      string // full path to the intermediate file
	GeneratedCode []byte // code which was run
	Output        []byte // resulting output (stdout + stderr)
}

type Framework interface {
	Info() *Info

	Initialize() (err error)

	SetModel(model *actr.Model) (err error)
	Model() (model *actr.Model)

	Run(initialBuffers InitialBuffers) (result *RunResult, err error)
	WriteModel(path string, initialBuffers InitialBuffers) (outputFileName string, err error)
}

type List map[string]Framework

// InitialBuffers is a map of buffer names to initial contents of the buffer.
// This is used when passing in user-defined initial contents e.g. through a web API.
type InitialBuffers map[string]string

// ParsedInitialBuffers is a map of buffer name to a parsed version of the initial contents.
// This is used when passing in user-defined initial contents e.g. through a web API.
type ParsedInitialBuffers map[string]*actr.Pattern

// Names returns all the names of the frameworks in the list.
func (l List) Names() (names []string) {
	names = make([]string, len(l))

	i := 0
	for k := range l {
		names[i] = k
		i++
	}

	return
}

// Exists checks if the framework is in the list.
func (l List) Exists(framework string) bool {
	for name := range l {
		if name == framework {
			return true
		}
	}

	return false
}

// IsValidFramework returns if the framework name is in our list of valid ones or not.
func IsValidFramework(frameworkName string) bool {
	return container.Contains(frameworkName, ValidFrameworks)
}
