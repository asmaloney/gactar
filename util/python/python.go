// Package python provides some helpers for dealing with running python & checking package availability.
package python

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
)

var (
	PythonRequiredVersion = 3

	ErrPython3NotFound = errors.New("python >= 3.0 not found. Please check your paths.")
)

type ErrPythonPackageNotFound struct {
	PackageName string
}

func (e *ErrPythonPackageNotFound) Error() string {
	return fmt.Sprintf("python package %q not found. Please ensure it is installed with pip or is in your PYTHONPATH env variable", e.PackageName)
}

func FindPython3(outputDetails bool) (path string, err error) {
	defer func() {
		if outputDetails && err == nil {
			fmt.Printf("> Found python: %q\n", path)

			output, execErr := executil.ExecCommand(path, "--version")
			if execErr != nil {
				return
			}

			fmt.Printf("> %s", output)
		}
	}()

	// See if it exists as "python3"
	path, err = filesystem.CheckForExecutable("python3")
	if err == nil {
		return
	}

	// See if it exists as "python"
	path, err = filesystem.CheckForExecutable("python")
	if err != nil {
		return
	}

	// We have "a python" now we need to check its version
	output, err := executil.ExecCommand(path, "--version")
	if err != nil {
		return
	}

	r := regexp.MustCompile(`(?i)^Python (\d+)\.\d+\.\d+`)
	match := r.FindStringSubmatch(output)
	if match == nil || len(match) < 2 {
		err = ErrPython3NotFound
		return
	}

	majorVersion, err := strconv.Atoi(match[1])
	if err != nil || majorVersion < PythonRequiredVersion {
		err = ErrPython3NotFound
		return
	}

	return
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
