// Package cmd implements our command-line commands for working with gactar.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/modes/defaultmode"
	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/executil"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/frameworkutil"
	"github.com/asmaloney/gactar/util/version"
)

var (
	ErrNoFrameworks = errors.New("could not create any frameworks - please check your installation")
	ErrNoInputFiles = errors.New("no input files specified on command line")
	ErrSilent       = errors.New("SilentErr")

	errNoFrameworks = errors.New("no frameworks specified on command line")

	validDebugOptions = []string{"lex", "parse", "exec"}

	flagEnv        = "./env"
	flagTemp       = ""
	flagFrameworks = []string{"all"}
	flagDebug      = []string{}
	flagNoColour   = false

	flagRun     = false
	flagVersion = false
)

type errRequiresSubcommand struct {
	command *cobra.Command
}

func (e errRequiresSubcommand) Error() string {
	message := fmt.Sprintf("%s command requires subcommand", e.command.Use)
	return chalk.ErrorBold(message)
}

type errVirtualEnvDoesNotExist struct {
	path string
}

func (e errVirtualEnvDoesNotExist) Error() string {
	message := fmt.Sprintf("virtual environment does not exist: %q", e.path)

	hint := fmt.Sprintf("hint: create and initialize a virtual environment using %q (use %q for details)",
		"gactar env setup", "gactar env setup --help")

	return chalk.ErrorBold(message) + "\n" + chalk.Default(hint)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "gactar",
	Short:         "A command-line tool for working with ACT-R models",
	Long:          "A proof-of-concept tool for creating ACT-R models using a declarative file format.",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flagNoColour {
			gchalk.SetLevel(gchalk.LevelNone)
		}

		amod.SetDebug(flagDebug)

		debugExec := slices.Contains(flagDebug, "exec")
		executil.SetDebug(debugExec)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if flagVersion {
			outputVersion()
			os.Exit(0)
		}

		settings, err := setupForRun(cmd)
		if err != nil {
			return err
		}

		s, err := defaultmode.Initialize(settings, args, flagRun)
		if err != nil {
			return err
		}

		err = s.Start()
		if err != nil {
			return err
		}

		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		if !errors.Is(err, ErrSilent) {
			chalk.PrintErr(err)
		}

		var requiresSubcommand errRequiresSubcommand
		if errors.As(err, &requiresSubcommand) {
			_ = requiresSubcommand.command.Help()
		}
		os.Exit(1)
	}
}

func init() {
	// We need to do some goofy stuff to handle formatting of flag vs. non-flag errors.
	// See: https://github.com/spf13/cobra/issues/914#issuecomment-548411337
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		chalk.PrintErr(err)
		cmd.Println(cmd.UsageString())
		return ErrSilent
	})

	// Global flags
	rootCmd.PersistentFlags().StringVar(&flagEnv, "env", flagEnv, "directory where ACT-R, pyactr, and other necessary files are installed")
	rootCmd.PersistentFlags().StringVar(&flagTemp, "temp", flagTemp, "directory for generated files (it will be created if it does not exist - defaults to <env>/gactar-temp)")
	rootCmd.PersistentFlags().StringSliceVarP(&flagFrameworks, "framework", "f", flagFrameworks,
		fmt.Sprintf("add framework - valid frameworks: %s", strings.Join(framework.ValidFrameworks, ", ")))
	rootCmd.PersistentFlags().StringSliceVarP(&flagDebug, "debug", "d", flagDebug,
		fmt.Sprintf("turn on debugging - valid options: %s", strings.Join(validDebugOptions, ", ")))
	rootCmd.PersistentFlags().BoolVar(&flagNoColour, "no-colour", false, "do not use colour output on command line")

	// Local flags - only run when this action is called directly.
	rootCmd.Flags().BoolVarP(&flagRun, "run", "r", false, "run the models after generating the code")
	rootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "output the version and quit")

	rootCmd.MarkFlagsMutuallyExclusive("run", "version")
	rootCmd.SetGlobalNormalizationFunc(normalizeAliasFlagsFunc)
}

func outputVersion() {
	version := fmt.Sprintf("gactar %s %s", "version", version.BuildVersion)
	fmt.Println(chalk.Bold(version))
}

// setupForRun sets up the virtual env, temp dir, and frameworks.
// It must be called by commands that are going to run gactar code (default, cli, & web).
func setupForRun(cmd *cobra.Command) (settings *cli.Settings, err error) {
	outputVersion()

	settings = &cli.Settings{
		Version: fmt.Sprintf("gactar %s %s", "version", version.BuildVersion),
	}

	envPath, err := setupVirtualEnvironment(cmd.Flags())
	if err != nil {
		return
	}

	settings.EnvPath = envPath

	// Create our temp dir. If it wasn't set, use <env>/gactar-temp.
	// createTempDirFromFlag() will expand our "temp" to an absolute path.
	tempPath, err := createTempDirFromFlag(cmd.Flags())
	if err != nil {
		return
	}

	settings.TempPath = tempPath

	frameworks, err := createFrameworks(settings, cmd.Flags())
	if err != nil {
		return
	}

	settings.ActiveFrameworks = frameworks

	return
}

func normalizeAliasFlagsFunc(flags *pflag.FlagSet, name string) pflag.NormalizedName {
	if name == "no-color" {
		name = "no-colour"
	}

	return pflag.NormalizedName(name)
}

// expandPathFlag expands the given path and sets it back in the flags.
func expandPathFlag(flags *pflag.FlagSet, name string) (path string, err error) {
	path, err = flags.GetString(name)
	if err != nil {
		return
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return
	}

	err = flags.Set(name, path)
	if err != nil {
		return
	}

	return
}

// createTempDirFromFlag looks up the "temp" flag in our flags, expands the path, and creates the dir.
func createTempDirFromFlag(flags *pflag.FlagSet) (path string, err error) {
	path, err = flags.GetString("temp")
	if err != nil {
		return
	}

	if path == "" {
		defaultTemp := fmt.Sprintf("%s%c%s", os.Getenv("VIRTUAL_ENV"), filepath.Separator, "gactar-temp")

		err = flags.Set("temp", defaultTemp)
		if err != nil {
			return
		}
	}

	path, err = expandPathFlag(flags, "temp")
	if err != nil {
		return
	}

	err = filesystem.CreateDir(path)
	if err != nil {
		return
	}

	return
}

// getVirtualEnvironmentPath gets the environment path from the flags and checks if it exists.
func getVirtualEnvironmentPath(flags *pflag.FlagSet) (envPath string, err error) {
	envPath, err = expandPathFlag(flags, "env")
	if err != nil {
		return "", err
	}

	if !filesystem.DirExists(envPath) {
		err = &errVirtualEnvDoesNotExist{path: envPath}
		return "", err
	}

	fmt.Print(chalk.Header("Using virtual environment: "))
	fmt.Printf("%q\n", envPath)

	return
}

// setupVirtualEnvironment will check that the environment path exists and set our PATH with it.
func setupVirtualEnvironment(flags *pflag.FlagSet) (path string, err error) {
	envPath, err := getVirtualEnvironmentPath(flags)
	if err != nil {
		return
	}

	err = cli.SetupPaths(envPath)
	if err != nil {
		return
	}

	path = envPath

	return
}

// createFrameworks will create the frameworks and return them as a list.
func createFrameworks(settings *cli.Settings, flags *pflag.FlagSet) (frameworks framework.List, err error) {
	list, err := flags.GetStringSlice("framework")
	if len(list) == 0 {
		err = errNoFrameworks
		return
	}

	// If the user asked for "all", then clear the list.
	// frameworkutil.CreateFrameworks() will create all valid ones.
	if slices.Contains(list, "all") {
		list = []string{}
	}

	frameworks = frameworkutil.CreateFrameworks(settings, list)

	if len(frameworks) == 0 {
		return framework.List{}, ErrNoFrameworks
	}

	return
}
