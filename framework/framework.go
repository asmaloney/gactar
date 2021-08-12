package framework

import (
	"fmt"
	"os/exec"

	"gitlab.com/asmaloney/gactar/actr"
)

type Framework interface {
	Initialize() (err error)
	SetModel(model *actr.Model) (err error)

	Run(initialGoal string) (output []byte, err error)
	WriteModel(path, initialGoal string) (outputFileName string, err error)
}

func CheckForExecutable(exe string) (path string, err error) {
	path, err = exec.LookPath(exe)
	if err != nil {
		err = fmt.Errorf("cannot find '%s' in your path", exe)
		return
	}

	return
}
