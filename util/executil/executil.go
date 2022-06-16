package executil

import (
	"fmt"
	"os/exec"
)

type ErrExecuteCommand struct {
	Output []byte
}

func (e ErrExecuteCommand) Error() string {
	return fmt.Sprintf("execution failed:\n%s", string(e.Output))
}

func ExecCommandWithCombinedOutput(name string, arg ...string) (output string, err error) {
	cmd := exec.Command(name, arg...)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		err = &ErrExecuteCommand{Output: outputBytes}
		return
	}

	output = string(outputBytes)
	return
}
