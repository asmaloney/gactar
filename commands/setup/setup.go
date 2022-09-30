package setup

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"

	"github.com/asmaloney/gactar/util/clicontext"
	"github.com/asmaloney/gactar/util/decompress"
	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/python"
)

var (
	ErrPathExists = errors.New("path already exists")
)

type ErrCCLSystem struct {
	OSName string
}

func (e ErrCCLSystem) Error() string {
	return fmt.Sprintf("no CCL compiler available for system %q", e.OSName)
}

type ErrExecuteCommand struct {
	Output []byte
}

func (e ErrExecuteCommand) Error() string {
	return fmt.Sprintf("execution failed:\n%s", string(e.Output))
}

func Setup(envPath string, dev bool) (err error) {
	fmt.Println("gactar Environment Setup\n---")
	fmt.Printf("Setting up an environment: %q\n", envPath)

	// Check if it already exists and error out
	if filesystem.DirExists(envPath) {
		err = fmt.Errorf("cannot set environment path to %q: %w", envPath, ErrPathExists)
		return
	} else {
		err = nil
	}

	// Create the virtual environment directory
	err = filesystem.CreateDir(envPath)
	if err != nil {
		return err
	}

	err = os.Chdir(envPath)
	if err != nil {
		return err
	}

	err = setupPython(envPath, dev)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
		// Don't return - we can still try to set up the Lisp compiler
	}

	err = setupLisp()
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	}

	return
}

func setupPython(envPath string, dev bool) (err error) {
	fmt.Println()
	fmt.Println("Setting up Python\n---")

	path, err := python.FindPython3(true)
	if err != nil {
		return
	}

	// Set up virtual environment
	fmt.Printf("> Setting up virtual environment: %q\n", envPath)
	_, err = executil.ExecCommand(path, "-m", "venv", envPath)
	if err != nil {
		return
	}

	err = clicontext.SetupPaths(envPath)
	if err != nil {
		return
	}

	fmt.Printf("> Reset PATH: %q\n", os.Getenv("PATH"))

	// Upgrade pip & install wheel
	var output string
	var errInstall error
	if runtime.GOOS == "windows" {
		// Windows fails on the pip upgrade for some reason, so leave it out
		fmt.Println("> Installing wheel...")
		output, errInstall = executil.ExecCommand("pip", "install", "wheel")

	} else {
		fmt.Println("> Upgrading pip & installing wheel...")
		output, errInstall = executil.ExecCommand("pip", "install", "--upgrade", "pip", "wheel")
	}
	if errInstall != nil {
		return errInstall
	}

	fmt.Print(output)

	// Install our requirements
	fmt.Println("> Installing pip packages...")

	requirementsFile := "requirements.txt"
	if dev {
		requirementsFile = "requirements-dev.txt"
	}

	output, err = executil.ExecCommand(
		"pip", "install", "-r",
		fmt.Sprintf("../install/%s", requirementsFile),
	)
	if err != nil {
		return
	}

	fmt.Print(output)

	return
}

func setupLisp() (err error) {
	fmt.Println()
	fmt.Println("Setting up Lisp\n---")

	// Download vanilla ACT-R
	repo := "github.com/asmaloney/ACT-R"
	version := "v7.27.0"
	archiveFile := fmt.Sprintf("actr-super-slim-%s.zip", version)
	urlStr := fmt.Sprintf("https://%s/releases/download/%s/%s", repo, version, archiveFile)
	actrURL, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	fmt.Printf("> Getting ACT-R %s from: %q\n", version, actrURL.String())

	err = filesystem.DownloadFile(actrURL, archiveFile)
	if err != nil {
		return
	}

	// Decompress ACT-R
	fmt.Println("> Unpacking ACT-R...")
	err = decompress.Unzip(archiveFile, "actr")
	if err != nil {
		return
	}

	// Download Clozure Common Lisp compiler (CCL)
	system := runtime.GOOS
	if system != "darwin" && system != "linux" && system != "windows" {
		return &ErrCCLSystem{OSName: system}
	}

	repo = "github.com/Clozure/ccl"
	extension := "tar.gz"
	version = "1.12.1"

	if system == "windows" {
		extension = "zip"
		version = "1.12" // version 1.12.1 is not compressed properly, so use older version
	}

	dirName := fmt.Sprintf("ccl-%s-%sx86", version, system)
	archiveFile = fmt.Sprintf("%s.%s", dirName, extension)
	urlStr = fmt.Sprintf("https://%s/releases/download/v%s/%s", repo, version, archiveFile)
	cclURL, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	fmt.Printf("> Getting Clozure Common Lisp (ccl) v%s for %s from: %q\n", version, system, cclURL.String())
	err = filesystem.DownloadFile(cclURL, archiveFile)
	if err != nil {
		return
	}

	// Decompress CCL
	fmt.Println("> Unpacking CCL...")
	if extension == "zip" {
		err = decompress.Unzip(archiveFile, "")
	} else {
		err = decompress.UntarFile(archiveFile, "")
	}
	if err != nil {
		return
	}

	return
}
