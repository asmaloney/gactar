package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/modes/defaultmode"
	"github.com/asmaloney/gactar/modes/shell"
	"github.com/asmaloney/gactar/modes/web"

	"github.com/asmaloney/gactar/util/clicontext"
	"github.com/asmaloney/gactar/util/container"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/frameworkutil"
	"github.com/asmaloney/gactar/util/version"
)

var (
	// "embed" cannot use relative paths, so we must declare this at the top level and pass into web.
	//go:embed examples/*.amod
	amodExamples embed.FS

	ErrNoFrameworks = errors.New("could not create any frameworks - please check your installation")
)

type ErrCmdLine struct {
	Message string
}

func (e *ErrCmdLine) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}

func main() {
	defaultPort := 8181
	defaultFramework := cli.NewStringSlice("all")

	app := &cli.App{
		Name:        "gactar",
		Usage:       "A command-line tool for working with ACT-R models",
		Description: "A proof-of-concept tool for creating ACT-R models using a declarative file format.",
		Version:     version.BuildVersion,
		Compiled:    time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Andy Maloney",
				Email: "asmaloney@gmail.com",
			},
		},
		Copyright:            "©2021 Andy Maloney",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "env", Value: "./env", Usage: "directory where ACT-R, pyactr, and other necessary files are installed", EnvVars: []string{"VIRTUAL_ENV"}},

			&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Usage: "turn on debugging output"},
			&cli.BoolFlag{Name: "ebnf", Usage: "output amod EBNF to stdout and quit"},
			&cli.PathFlag{Name: "temp", Value: "./gactar-temp", Usage: "directory for generated files (it will be created if it does not exist)"},

			&cli.StringSliceFlag{
				Name:    "framework",
				Aliases: []string{"f"},
				Value:   defaultFramework,
				Usage:   fmt.Sprintf("add framework - valid frameworks: %s", strings.Join(framework.ValidFrameworks, ", ")),
			},

			// CLI mode
			&cli.BoolFlag{Name: "run", Aliases: []string{"r"}, Category: "Mode: CLI", Usage: "run the models after generating the code"},

			// CLI (interactive) mode
			&cli.BoolFlag{Name: "interactive", Aliases: []string{"i"}, Category: "Mode: CLI (interactive)", Usage: "run an interactive shell"},

			// Web mode
			&cli.BoolFlag{Name: "web", Aliases: []string{"w"}, Category: "Mode: Web", Usage: "start a web server to run in a browser"},
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Category: "Mode: Web", Value: defaultPort, Usage: "port to run the web server on"},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				amod.SetDebug(true)
			}

			if c.Bool("ebnf") {
				amod.OutputEBNF()
				return nil
			}

			if c.Bool("web") && c.Bool("interactive") {
				err := &ErrCmdLine{Message: "cannot run 'web' and 'interactive' at the same time"}
				return cli.Exit(err.Error(), 1)
			}

			err := setupVirtualEnvironment(c)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			// Create our temp dir. This will expand our "temp" to an absolute path.
			err = clicontext.CreateTempDir(c)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			frameworks, err := createFrameworks(c)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if c.Bool("web") {
				err = handleWeb(c, frameworks)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
				return nil
			}

			if c.Int("port") != defaultPort {
				fmt.Println("info: --port only applies when using --web")
			}

			if c.Bool("interactive") {
				err = handleInteractive(c, frameworks)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}

				return nil
			}

			// We are not interactive or web, so simply generate the output files.
			err = handleDefault(c, frameworks)
			if err != nil {
				clErr := &ErrCmdLine{Message: err.Error()}
				return cli.Exit(clErr.Error(), 1)
			}

			return nil
		},
	}

	// Used to output command line options for documentation.
	// fmt.Println(app.ToMarkdown())

	app.Run(os.Args) //nolint - exits are handled with cli.Exit()
}

// setupVirtualEnvironment will set our paths to our virtual environment path.
func setupVirtualEnvironment(ctx *cli.Context) (err error) {
	envPath, err := clicontext.ExpandPath(ctx, "env")
	if err != nil {
		return
	}

	if !filesystem.DirExists(envPath) {
		err = &ErrCmdLine{Message: "virtual environment does not exist"}
		err = fmt.Errorf("%w: %q", err, envPath)
		return
	}

	fmt.Printf("Using virtual environment: %q\n", envPath)

	binPath := filepath.Join(envPath, "bin")
	os.Setenv("PATH", fmt.Sprintf("%s%c%s", binPath, os.PathListSeparator, os.Getenv("PATH")))
	os.Setenv("VIRTUAL_ENV", envPath)

	return
}

func createFrameworks(cli *cli.Context) (frameworks framework.List, err error) {
	list := cli.StringSlice("framework")
	if len(list) == 0 {
		err = &ErrCmdLine{Message: "no frameworks specified on command line"}
		return
	}

	// If the user asked for "all", then clear the list.
	// frameworkutil.CreateFrameworks() will create all valid ones.
	if container.Contains("all", list) {
		list = []string{}
	}

	frameworks = frameworkutil.CreateFrameworks(cli, list)

	if len(frameworks) == 0 {
		return framework.List{}, ErrNoFrameworks
	}

	return
}

func handleWeb(ctx *cli.Context, frameworks framework.List) (err error) {
	w, err := web.Initialize(ctx, frameworks, &amodExamples)
	if err != nil {
		return err
	}

	err = w.Start()
	if err != nil {
		return err
	}

	return
}

func handleInteractive(ctx *cli.Context, frameworks framework.List) (err error) {
	s, err := shell.Initialize(ctx, frameworks)
	if err != nil {
		return err
	}

	err = s.Start()
	if err != nil {
		return err
	}

	return
}

func handleDefault(ctx *cli.Context, frameworks framework.List) (err error) {
	s, err := defaultmode.Initialize(ctx, frameworks)
	if err != nil {
		return err
	}

	err = s.Start()
	if err != nil {
		return err
	}

	return
}
