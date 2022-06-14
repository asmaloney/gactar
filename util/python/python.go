// Package python provides some helpers for dealing with running python & checking package availability.
package python

import (
	"fmt"
	"os/exec"
)

type ErrPythonPackageNotFound struct {
	PackageName string
}

func (e *ErrPythonPackageNotFound) Error() string {
	return fmt.Sprintf("python package %q not found. Please ensure it is installed with pip or is in your PYTHONPATH env variable", e.PackageName)
}

// CheckForPackage checks for the proper installation of the named package.
func CheckForPackage(python, packageName string) (err error) {
	importCmd := fmt.Sprintf("import %s", packageName)

	cmd := exec.Command(python, "-c", importCmd)

	err = cmd.Run()
	if err != nil {
		err = &ErrPythonPackageNotFound{PackageName: packageName}
		return
	}

	return
}
