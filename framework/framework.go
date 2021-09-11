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
	SetModel(model *actr.Model) (err error)

	Run(initialGoal string) (generatedCode, output []byte, err error)
	WriteModel(path, initialGoal string) (outputFileName string, err error)
}

type List map[string]Framework

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
