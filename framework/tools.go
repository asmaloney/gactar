package framework

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gitlab.com/asmaloney/gactar/actr"
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

// Float64Str takes a float and returns a string of the minimal representation.
// e.g. 2.5000 becomes "2.5"
func Float64Str(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
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

func PythonValuesToStrings(values *[]*actr.Value, quoteStrings bool) []string {
	str := make([]string, len(*values))
	for i, v := range *values {
		if v.Var != nil {
			str[i] = strings.TrimPrefix(*v.Var, "?")
		} else if v.Str != nil {
			if quoteStrings {
				str[i] = fmt.Sprintf("'%s'", *v.Str)
			} else {
				str[i] = *v.Str
			}
		} else if v.Number != nil {
			str[i] = *v.Number
		}
		// v.ID should not be possible because of validation
	}

	return str
}
