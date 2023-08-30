package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/decompress"
	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/python"
)

const (
	// JSON file used to override our defaults
	SUPPORT_TOOLS_FILE = "install/support-tools.json"

	// Default ACT-R repo
	ACTR_REPO = "github.com/asmaloney/ACT-R"

	// Default ACT-R release version
	ACTR_VERSION = "7.27.7"

	// Default Clozure Common Lisp repo
	CCL_REPO = "github.com/Clozure/ccl"

	// Default Clozure Common Lisp release version
	CCL_VERSION = "1.12.2"
)

var (
	errPathExists = errors.New("path already exists")

	flagSetupDev = false

	flagUpdateAll = false

	flagUpdatePython         = false
	flagUpdatePythonPackages = false
	flagUpdateDev            = false

	// names we allow in the support-tools file
	validTools = []string{"ACT-R", "CCL"}

	defaultToolInfo = make(toolInfoMap, len(validTools))
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Setup & maintain an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = errRequiresSubcommand{command: cmd}
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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an environment",
	PreRun: func(cmd *cobra.Command, args []string) {
		devSet, _ := cmd.Flags().GetBool("dev")
		allSet, _ := cmd.Flags().GetBool("all")
		if devSet && !allSet {
			cmd.MarkFlagRequired("pip")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		envPath, err := getVirtualEnvironmentPath(cmd.Flags())
		if err != nil {
			chalk.PrintErr(err)
			return
		}

		err = runUpdate(envPath)
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

// toolInfoList is used to read tool version info from a file
type toolInfoList struct {
	List []toolInfo `json:"tool-info"`
}

// toolInfoMap maps a tool name to its info
type toolInfoMap map[string]toolInfo

// toolInfo stores version info about a tool
type toolInfo struct {
	Name    string `json:"name"`
	RepoURL string `json:"repo-url"`
	Version string `json:"version"`
}

func init() {
	setDefaultToolInfo()

	setupCmd.Flags().BoolVar(&flagSetupDev, "dev", false, "install dev packages")
	setupCmd.Flags().StringP("path", "p", "./env", "directory for env files (it will be created if it does not exist)")

	envCmd.AddCommand(setupCmd)

	updateCmd.Flags().BoolVar(&flagUpdateAll, "all", false, "update all tools & packages")
	updateCmd.Flags().BoolVar(&flagUpdatePython, "python", false, "update python version")
	updateCmd.Flags().BoolVar(&flagUpdatePythonPackages, "pip", false, "update python packages")
	updateCmd.Flags().BoolVar(&flagUpdateDev, "dev", false, "update dev packages")

	envCmd.AddCommand(updateCmd)

	rootCmd.AddCommand(envCmd)
}

// setDefaultToolInfo fills in the default tool versions as a fallback in case the external
// file doesn't exist or fails to parse.
func setDefaultToolInfo() {
	defaultToolInfo["ACT-R"] = toolInfo{
		Name:    "ACT-R",
		RepoURL: ACTR_REPO,
		Version: ACTR_VERSION,
	}

	defaultToolInfo["CCL"] = toolInfo{
		Name:    "CCL",
		RepoURL: CCL_REPO,
		Version: CCL_VERSION,
	}
}

func readToolInfo() (toolInfo toolInfoMap, err error) {
	file, err := os.ReadFile(SUPPORT_TOOLS_FILE)
	if err != nil {
		return
	}

	data := toolInfoList{}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return
	}

	if len(data.List) == 0 {
		chalk.PrintWarningStr(fmt.Sprintf("failed to find version information in %q; using defaults", SUPPORT_TOOLS_FILE))
		return defaultToolInfo, err
	}

	toolInfo = defaultToolInfo

	// Get our version from the support-tools file
	for _, info := range data.List {
		if !slices.Contains(validTools, info.Name) {
			chalk.PrintWarningStr(fmt.Sprintf("invalid tool name found in support-tools file: %q", info.Name))
			continue
		}

		tool := toolInfo[info.Name]
		tool.Version = info.Version
		toolInfo[info.Name] = tool
	}

	return
}

func runSetup(envPath string, dev bool) (err error) {
	fmt.Println("gactar Environment Setup\n---")
	fmt.Printf("Setting up an environment: %q\n", envPath)

	// Check if it already exists and error out
	if filesystem.DirExists(envPath) {
		err = fmt.Errorf("cannot set environment path to %q: %w", envPath, errPathExists)
		return
	}

	tools, err := readToolInfo()
	if err != nil {
		return err
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
		chalk.PrintErr(err)
		// Don't return - we can still try to set up the Lisp compiler
	}

	err = setupLisp(tools)
	if err != nil {
		return err
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

func setupLisp(tools toolInfoMap) (err error) {
	fmt.Println()
	fmt.Println("Setting up Lisp\n---")

	// Download vanilla ACT-R
	actrInfo := tools["ACT-R"]
	archiveFile := fmt.Sprintf("actr-super-slim-v%s.zip", actrInfo.Version)

	err = downloadGitHubRelease("ACT-R", actrInfo.RepoURL, actrInfo.Version, archiveFile, "actr")
	if err != nil {
		return
	}

	// Download Clozure Common Lisp compiler (CCL)
	cclInfo := tools["CCL"]
	system := runtime.GOOS
	if system != "darwin" && system != "linux" && system != "windows" {
		return &errCCLSystem{OSName: system}
	}

	extension := "tar.gz"
	if system == "windows" {
		extension = "zip"
	}

	dirName := fmt.Sprintf("ccl-%s-%sx86", cclInfo.Version, system)
	archiveFile = fmt.Sprintf("%s.%s", dirName, extension)

	err = downloadGitHubRelease("Clozure Common Lisp (ccl)", cclInfo.RepoURL, cclInfo.Version, archiveFile, "")
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

func runUpdate(envPath string) (err error) {
	fmt.Println("gactar Environment Update\n---")
	fmt.Printf("Updating environment: %q\n", envPath)

	if flagUpdateAll || flagUpdatePython {
		err = updatePython(envPath)
		if err != nil {
			return
		}
	}

	if flagUpdateAll || flagUpdatePythonPackages {
		err = updatePipPackages(envPath)
		if err != nil {
			return
		}
	}

	return
}

func updatePython(envPath string) (err error) {
	fmt.Println()
	fmt.Println("Updating Python\n---")

	path, err := python.FindPython3(true)
	if err != nil {
		return
	}

	output, err := executil.ExecCommand(path, "-m", "venv", "--upgrade", envPath)
	if err != nil {
		return
	}

	fmt.Print(output)

	return
}

func updatePipPackages(envPath string) (err error) {
	fmt.Println()
	fmt.Println("Updating Python pip packages\n---")

	err = cli.SetupPaths(envPath)
	if err != nil {
		return
	}

	file := "install/requirements.txt"
	if flagUpdateDev {
		file = "install/requirements-dev.txt"
	}

	output, err := executil.ExecCommand("pip", "install", "-U", "-r", file)
	if err != nil {
		return
	}

	fmt.Print(output)

	return
}
