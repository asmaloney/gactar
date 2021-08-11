package framework

import (
	"fmt"
	"os/exec"
	"strings"
)

// Some tools for working with Python

// PythonIdentify outputs version info and the path to the python executable.
func PythonIdentify() {
	cmd := exec.Command("python3", "--version")
	output, _ := cmd.CombinedOutput()

	version := strings.TrimSpace(string(output))

	cmd = exec.Command("which", "python3")
	output, _ = cmd.CombinedOutput()

	fmt.Printf("Using %s from %s", version, string(output))
}

// PythonCheckForPackage checks for the proper installation of the named package.
func PythonCheckForPackage(packageName string) (err error) {
	importCmd := fmt.Sprintf("import %s", packageName)

	cmd := exec.Command("python3", "-c", importCmd)

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("python package '%s' not found. Please ensure it is installed with pip or is in your PYTHONPATH env variable", packageName)
		return
	}

	return
}
