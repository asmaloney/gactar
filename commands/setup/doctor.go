package setup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/lisp"
	"github.com/asmaloney/gactar/util/python"
	"github.com/asmaloney/gactar/util/version"
)

var (
	ErrHealthCheckFailed = errors.New("gactar health check failed")
)

func Doctor(envPath string) (err error) {
	fmt.Println("gactar Environment Doctor\n-----")
	fmt.Printf("gactar %s\n", version.BuildVersion)
	fmt.Printf("Checking %q for problems...\n", envPath)
	fmt.Printf("PATH: %q\n", os.Getenv("PATH"))

	if !filesystem.DirExists(envPath) {
		return &filesystem.ErrDirDoesNotExist{DirName: envPath}
	}

	e := os.Chdir(envPath)
	if e != nil {
		fmt.Println(e.Error())
		return ErrHealthCheckFailed
	}

	pythonPath, e := checkPython()
	if e != nil {
		fmt.Println(e.Error())
		err = ErrHealthCheckFailed
	}

	e = checkCLL()
	if e != nil {
		fmt.Println(e.Error())
		err = ErrHealthCheckFailed
	}

	e = checkFrameworks(envPath, pythonPath)
	if e != nil {
		err = ErrHealthCheckFailed
	}

	return
}

func checkPython() (path string, err error) {
	fmt.Println()
	fmt.Println("Checking Python\n---")

	path, err = python.FindPython3(true)
	if err != nil {
		return
	}

	return
}

func checkCLL() (err error) {
	fmt.Println()
	fmt.Println("Checking Clozure Common Lisp (ccl) compiler\n---")

	cclExecutableName, err := lisp.GetExecutableName()
	if err != nil {
		return
	}

	fmt.Printf("> Looking for %q...\n", cclExecutableName)

	exePath, err := filesystem.CheckForExecutable(cclExecutableName)
	if err != nil {
		return
	}

	fmt.Printf("> Found ccl: %q\n", exePath)

	output, err := executil.ExecCommand(exePath, "--version")
	if err != nil {
		return err
	}
	fmt.Printf("> %s", output)

	return
}

func checkFrameworks(envPath, pythonPath string) (err error) {
	fmt.Println()
	fmt.Println("Checking Frameworks\n---")

	atLeastOne := false
	fmt.Println("> Checking for python_actr (ccm) package...")
	e := python.CheckForPackage(pythonPath, "python_actr")
	if e != nil {
		fmt.Println(e.Error())
		fmt.Println("> NOTE: python_actr (ccm) not available")
	} else {
		atLeastOne = true
		fmt.Println("> ...found")
	}

	fmt.Println("> Checking for pyactr package...")
	e = python.CheckForPackage(pythonPath, "pyactr")
	if e != nil {
		fmt.Println(e.Error())
		fmt.Println("> NOTE: pyactr not available")
	} else {
		atLeastOne = true
		fmt.Println("> ...found")
	}

	fmt.Println("> Checking for ACT-R source...")
	actrDir := filepath.Join(envPath, "actr")
	if !filesystem.DirExists(actrDir) {
		e = &filesystem.ErrDirDoesNotExist{DirName: actrDir}
		fmt.Println(e.Error())
		fmt.Println("> NOTE: vanilla ACT-R not available")
	} else {
		atLeastOne = true
		fmt.Printf("> ...found: %q\n", actrDir)
	}

	if !atLeastOne {
		fmt.Println("> ERROR: could not find any frameworks")
		err = ErrHealthCheckFailed
	}

	return
}
