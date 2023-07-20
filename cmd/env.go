package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/decompress"
	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/python"
)

const (
	// ACT-R release version from https://github.com/asmaloney/ACT-R
	ACTR_VERSION = "7.27.7"

	// Clozure Common Lisp release version from https://github.com/Clozure/ccl
	CCL_VERSION     = "1.12.1"
	CCL_VERSION_WIN = "1.12" // CCL version 1.12.1 is not compressed properly, so try an older version
)

var (
	errPathExists = errors.New("path already exists")

	flagSetupDev = false
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Setup & maintain an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = errRequiresSubcommand{command: "env"}
		chalk.PrintErr(err)

		return
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup an environment",
	Run: func(cmd *cobra.Command, args []string) {
		envPath, err := expandPathFlag(cmd.Flags(), "path")
		if err != nil {
			chalk.PrintErr(err)
			return
		}

		err = runSetup(envPath, flagSetupDev)
		if err != nil {
			chalk.PrintErr(err)
			return
		}
	},
}

type errCCLSystem struct {
	OSName string
}

func (e errCCLSystem) Error() string {
	return fmt.Sprintf("no CCL compiler available for system %q", e.OSName)
}

func init() {
	setupCmd.Flags().StringP("path", "p", "./env", "directory for env files (it will be created if it does not exist)")
	setupCmd.Flags().BoolVar(&flagSetupDev, "dev", false, "install dev packages")

	envCmd.AddCommand(setupCmd)

	rootCmd.AddCommand(envCmd)
}

func runSetup(envPath string, dev bool) (err error) {
	fmt.Println("gactar Environment Setup\n---")
	fmt.Printf("Setting up an environment: %q\n", envPath)

	// Check if it already exists and error out
	if filesystem.DirExists(envPath) {
		err = fmt.Errorf("cannot set environment path to %q: %w", envPath, errPathExists)
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

	err = cli.SetupPaths(envPath)
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
	extension := "zip"
	archiveFile := fmt.Sprintf("actr-super-slim-v%s.%s", ACTR_VERSION, extension)

	err = downloadGitHubRelease("ACT-R", repo, ACTR_VERSION, archiveFile, "actr")
	if err != nil {
		return
	}

	// Download Clozure Common Lisp compiler (CCL)
	system := runtime.GOOS
	if system != "darwin" && system != "linux" && system != "windows" {
		return &errCCLSystem{OSName: system}
	}

	repo = "github.com/Clozure/ccl"
	extension = "tar.gz"
	version := CCL_VERSION

	if system == "windows" {
		extension = "zip"
		version = CCL_VERSION_WIN
	}

	dirName := fmt.Sprintf("ccl-%s-%sx86", version, system)
	archiveFile = fmt.Sprintf("%s.%s", dirName, extension)

	err = downloadGitHubRelease("Clozure Common Lisp (ccl)", repo, version, archiveFile, "")
	if err != nil {
		return
	}

	return
}

func downloadGitHubRelease(name, repo, version, archiveFile, target string) (err error) {
	urlStr := fmt.Sprintf("https://%s/releases/download/v%s/%s", repo, version, archiveFile)
	url, err := url.Parse(urlStr)
	if err != nil {
		return
	}

	fmt.Printf("> Getting %s v%s from: %q\n", name, version, url.String())

	err = filesystem.DownloadFile(url, archiveFile)
	if err != nil {
		return
	}

	// Decompress
	fmt.Printf("> Unpacking %s...\n", name)
	extension := filepath.Ext(archiveFile)
	if extension == ".zip" {
		err = decompress.Unzip(archiveFile, target)
	} else {
		err = decompress.UntarFile(archiveFile, target)
	}
	if err != nil {
		return
	}

	return
}
