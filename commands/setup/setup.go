package setup

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

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

func Setup(envPath string) (err error) {
	fmt.Println("Setup an environment: ", envPath)

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

	binPath := filepath.Join(envPath, "bin")
	os.Setenv("PATH", fmt.Sprintf("%s%c%s", binPath, os.PathListSeparator, os.Getenv("PATH")))
	if err != nil {
		return err
	}

	err = os.Chdir(envPath)
	if err != nil {
		return err
	}

	err = setupPython(envPath)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
		// Don't return - we can still try to set up the Lisp compiler
	}

	err = setupLisp(envPath)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	}

	return
}

func setupPython(envPath string) (err error) {
	fmt.Println()
	fmt.Println("Setting up Python\n---")

	// Run some python 3
	path, err := python.FindPython3()
	if err != nil {
		return
	}

	fmt.Printf("> Found python3: %q\n", path)
	output, err := executil.ExecCommandWithCombinedOutput(path, "--version")
	if err != nil {
		return
	}

	fmt.Print(output)

	// Set up virtual environment
	fmt.Printf("> Setting up virtual environment: %q\n", path)
	_, err = executil.ExecCommandWithCombinedOutput(path, "-m", "venv", envPath)
	if err != nil {
		return
	}

	os.Setenv("VIRTUAL_ENV", envPath)

	// Upgrade pip
	fmt.Println("> Upgrading pip...")
	output, err = executil.ExecCommandWithCombinedOutput("pip", "install", "--upgrade", "pip", "wheel")
	if err != nil {
		return
	}

	fmt.Print(output)

	// Install our requirements
	fmt.Println("> Installing pip packages...")
	output, err = executil.ExecCommandWithCombinedOutput("pip", "install", "-r", "../install/requirements.txt")
	if err != nil {
		return
	}

	fmt.Print(output)

	return
}

func setupLisp(envPath string) (err error) {
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
