// Package executil implements a function for executing a command line.
package executil

import (
	"fmt"
	"os/exec"
)

type ErrExecuteCommand struct {
	Output string
}

func (e ErrExecuteCommand) Error() string {
	return fmt.Sprintf("execution failed:\n%s", e.Output)
}

func ExecCommand(name string, arg ...string) (output string, err error) {
	cmd := exec.Command(name, arg...)
	outputBytes, err := cmd.CombinedOutput()
	output = string(outputBytes)
	if err != nil {
		err = &ErrExecuteCommand{Output: output}
		return
	}

	return
}
