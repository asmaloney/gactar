package framework

import (
	"fmt"
	"os/exec"

	"github.com/asmaloney/gactar/actr"
)

// ValidFrameworks lists the valid options for choosing frameworks on the command line and in the
// interactive case.
var ValidFrameworks = []string{"all", "ccm", "pyactr", "vanilla"}

type Framework interface {
	Initialize() (err error)
	Name() string

	SetModel(model *actr.Model) (err error)
	Model() (model *actr.Model)

	Run(initialBuffers InitialBuffers) (generatedCode, output []byte, err error)
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

func CheckForExecutable(exe string) (path string, err error) {
	path, err = exec.LookPath(exe)
	if err != nil {
		err = fmt.Errorf("cannot find '%s' in your path", exe)
		return
	}

	return
}

func IsValidFramework(framework string) bool {
	for _, f := range ValidFrameworks {
		if f == framework {
			return true
		}
	}

	return false
}
