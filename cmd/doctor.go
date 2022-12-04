package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/lisp"
	"github.com/asmaloney/gactar/util/python"
	"github.com/asmaloney/gactar/util/version"
)

var (
	ErrHealthCheckFailed = errors.New("gactar health check failed")
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check an environment for problems",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		envPath, err := expandPathFlag(cmd.Flags(), "path")
		if err != nil {
			return
		}

		err = cli.SetupPaths(envPath)
		if err != nil {
			return
		}

		err = runDoctor(envPath)
		fmt.Println()
		if err != nil {
			return
		}

		fmt.Print(chalk.Success("Health check passed."))
		fmt.Println(" Go forth and model!")

		return
	},
}

func init() {
	envCmd.AddCommand(doctorCmd)

	doctorCmd.Flags().StringP("path", "p", "./env", "environment to check")
}

func outputSectionHeader(text string) {
	fmt.Println(chalk.BlueUnderline(text))
	if !chalk.HasColor() {
		fmt.Println("---")
	}
}

func runDoctor(envPath string) (err error) {
	fmt.Println(chalk.BlueBoldUnderline("gactar Environment Doctor"))

	if !chalk.HasColor() {
		fmt.Println("-----")
	}

	fmt.Print(chalk.Header("gactar version: "))
	fmt.Println(version.BuildVersion)

	fmt.Print(chalk.Header("Shell PATH: "))
	fmt.Println(os.Getenv("PATH"))

	fmt.Print(chalk.Header("Environment: "))
	fmt.Println(envPath)

	if !filesystem.DirExists(envPath) {
		return &filesystem.ErrDirDoesNotExist{DirName: envPath}
	}

	e := os.Chdir(envPath)
	if e != nil {
		chalk.PrintErrLight(e)
		return ErrHealthCheckFailed
	}

	pythonPath, e := checkPython()
	if e != nil {
		chalk.PrintErrLight(e)
		err = ErrHealthCheckFailed
	}

	e = checkCLL()
	if e != nil {
		chalk.PrintErrLight(e)
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
	outputSectionHeader("Python")

	fmt.Printf("> Looking for %s...\n", chalk.Italic("python 3.x"))

	path, err = python.FindPython3(true)
	if err != nil {
		return
	}

	return
}

func checkCLL() (err error) {
	fmt.Println()
	outputSectionHeader("Clozure Common Lisp (ccl) compiler")

	cclExecutableName, err := lisp.GetExecutableName()
	if err != nil {
		return
	}

	fmt.Printf("> Looking for %s...\n", chalk.Italic(cclExecutableName))

	exePath, err := filesystem.CheckForExecutable(cclExecutableName)
	if err != nil {
		return
	}

	fmt.Printf("> Found ccl: %s\n", exePath)

	output, err := executil.ExecCommand(exePath, "--version")
	if err != nil {
		return err
	}
	fmt.Printf(">   %s", output)

	return
}

func checkForPythonPackage(packageName, pythonPath string) bool {
	fmt.Printf("> Checking for %s package...\n", chalk.Italic(packageName))
	err := python.CheckForPackage(pythonPath, packageName)
	if err != nil {
		chalk.PrintErrLight(err)
		chalk.PrintWarningStr(fmt.Sprintf("> NOTE: %s not available\n", packageName))
		return false
	}

	fmt.Println(">   ...found")
	return true
}

func checkFrameworks(envPath, pythonPath string) (err error) {
	fmt.Println()
	outputSectionHeader("Frameworks")

	atLeastOne := checkForPythonPackage("python_actr", pythonPath)
	atLeastOne = checkForPythonPackage("pyactr", pythonPath) || atLeastOne

	fmt.Printf("> Checking for %s source...\n", chalk.Italic("ACT-R"))
	actrDir := filepath.Join(envPath, "actr")
	if !filesystem.DirExists(actrDir) {
		e := &filesystem.ErrDirDoesNotExist{DirName: actrDir}
		chalk.PrintErrLight(e)
		chalk.PrintWarningStr("> NOTE: vanilla ACT-R not available")
	} else {
		atLeastOne = true
		fmt.Printf(">   ...found: %s\n", actrDir)
	}

	if !atLeastOne {
		chalk.PrintErrStr("could not find any frameworks")
		err = ErrHealthCheckFailed
	}

	return
}
