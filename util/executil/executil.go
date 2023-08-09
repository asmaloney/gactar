// Package executil implements a function for executing a command line.
package executil

import (
	"fmt"
	"os/exec"
)

var (
	debugging bool = false
)

type ErrExecuteCommand struct {
	Output string
}

func (e ErrExecuteCommand) Error() string {
	return fmt.Sprintf("execution failed:\n%s", e.Output)
}

func ExecCommand(name string, arg ...string) (output string, err error) {
	cmd := exec.Command(name, arg...)

	if debugging {
		fmt.Printf("Executing: %s\n", cmd.String())
	}

	outputBytes, err := cmd.CombinedOutput()
	output = string(outputBytes)

	if debugging {
		if err != nil {
			fmt.Printf("Exec FAIL: %s\n%s\n", cmd.String(), output)
		} else {
			fmt.Printf("Exec SUCCESS: %s\n", cmd.String())
		}
	}

	if err != nil {
		err = &ErrExecuteCommand{Output: output}
		return
	}

	return
}

// SetDebug turns debugging on and off. This will output the command before executing it.
func SetDebug(debug bool) {
	debugging = debug
}
