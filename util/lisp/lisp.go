package lisp

import (
	"fmt"
	"runtime"
)

type ErrCCLSystem struct {
	OSName string
	OSArch string
}

func (e ErrCCLSystem) Error() string {
	return fmt.Sprintf("no CCL compiler available for system %s/%s", e.OSName, e.OSArch)
}

func GetExecutableName() (exeName string, err error) {
	// We only support 64-bit. Nobody still uses 32-bit, right?
	if runtime.GOARCH == "386" {
		err = &ErrCCLSystem{OSName: runtime.GOOS, OSArch: runtime.GOARCH}
		return
	}

	// These executable names come from the CCL release tarballs.
	switch runtime.GOOS {
	case "darwin":
		exeName = "dx86cl64"

	case "linux":
		exeName = "lx86cl64"

	case "windows":
		exeName = "wx86cl64.exe"

	default:
		err = &ErrCCLSystem{OSName: runtime.GOOS, OSArch: runtime.GOARCH}
		return
	}

	return
}
