package framework

import (
	"fmt"
	"os/exec"
	"strings"
)

// Some tools for working with our frameworks

// IdentifyYourself outputs version info and the path to an executable.
func IdentifyYourself(frameworkName, exeName string) {
	cmd := exec.Command(exeName, "--version")
	output, _ := cmd.CombinedOutput()

	version := strings.TrimSpace(string(output))

	cmd = exec.Command("which", exeName)
	output, _ = cmd.CombinedOutput()

	fmt.Printf("%s: Using %s from %s", frameworkName, version, string(output))
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
