package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/framework/ccm_pyactr"
	"github.com/asmaloney/gactar/framework/pyactr"
	"github.com/asmaloney/gactar/framework/vanilla_actr"
	"github.com/asmaloney/gactar/shell"
	"github.com/asmaloney/gactar/version"
	"github.com/asmaloney/gactar/web"

	"github.com/asmaloney/gactar/util/clicontext"
	"github.com/asmaloney/gactar/util/container"
	"github.com/asmaloney/gactar/util/validate"
)

// "embed" cannot use relative paths, so we must declare this at the top level and pass into web.
//go:embed examples/*.amod
var amodExamples embed.FS

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
			&cli.StringFlag{Name: "env", Usage: "directory where ACT-R, pyactr, and other necessary files are installed", EnvVars: []string{"VIRTUAL_ENV"}},

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
				err := errors.New("cannot run 'web' and 'interactive' at the same time")
				fmt.Println(err.Error())
				return err
			}

			// Create our temp dir. This will expand our "temp" to an absolute path.
			err := clicontext.CreateTempDir(c)
			if err != nil {
				return err
			}

			frameworks, err := createFrameworks(c)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			if c.Bool("web") {
				err := handleWeb(c, frameworks)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
			}

			if c.Int("port") != defaultPort {
				fmt.Println("info: --port only applies when using --web")
			}

			if c.Bool("interactive") {
				err := handleInteractive(c, frameworks)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				return err
			}

			// We are not interactive or web, so simply generate the output files.
			err = handleDefault(c, frameworks)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			return nil
		},
	}

	// Used to output command line options for documentation.
	// fmt.Println(app.ToMarkdown())

	app.Run(os.Args)
}

func createFrameworks(cli *cli.Context) (frameworks framework.List, err error) {
	list := cli.StringSlice("framework")
	if len(list) == 0 {
		err = fmt.Errorf("no frameworks specified on command line")
		return
	}

	list = container.UniqueAndSorted(list)

	if list[0] == "all" {
		list = framework.ValidNamedFrameworks()
	}

	frameworks = make(framework.List, len(list))

	for _, f := range list {
		var createErr error
		switch f {
		case "ccm":
			frameworks["ccm"], createErr = ccm_pyactr.New(cli)
		case "pyactr":
			frameworks["pyactr"], createErr = pyactr.New(cli)
		case "vanilla":
			frameworks["vanilla"], createErr = vanilla_actr.New(cli)
		default:
			err = fmt.Errorf("unknown framework: %s", f)
			return framework.List{}, err
		}

		if createErr != nil {
			fmt.Println(createErr.Error())
		}
	}

	if len(frameworks) == 0 {
		err = fmt.Errorf("could not create any frameworks - please check your installation")
		return framework.List{}, err
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
	cli.ShowVersion(ctx)

	// Check if files exist first
	files := ctx.Args().Slice()

	if len(files) == 0 {
		err = fmt.Errorf("error: no input files specified on command line")
		return
	}

	existingFiles := files[:0]
	for _, file := range files {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("error: file does not exist - %q\n", file)
			continue
		}

		existingFiles = append(existingFiles, file)
	}

	if len(existingFiles) == 0 {
		err = fmt.Errorf("error: no files to process")
		return
	}

	tempPath := ctx.Path("temp")
	fmt.Printf("Intermediate file path: %q\n", tempPath)

	err = generateCode(frameworks, existingFiles, tempPath, ctx.Bool("run"))
	if err != nil {
		return err
	}

	if ctx.Bool("run") {
		runCode(frameworks)
	}
	return
}

func generateCode(frameworks framework.List, files []string, outputDir string, runCode bool) (err error) {
	modelMap := map[string]*actr.Model{}

	for _, file := range files {
		fmt.Printf("Generating model for %s\n", file)
		model, log, err := amod.GenerateModelFromFile(file)
		if err != nil {
			fmt.Print(log)
			continue
		}

		// When using "-r" the goal must be initialized in the code.
		validate.Goal(model, "", log)

		fmt.Print(log)

		modelMap[file] = model
	}

	if len(modelMap) == 0 {
		err = errors.New("no valid models to run")
		return
	}

	for _, f := range frameworks {
		err = f.Initialize()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for file, model := range modelMap {
			fmt.Printf("\t- generating code for %s\n", file)

			log := f.ValidateModel(model)
			fmt.Print(log)
			if log.HasError() {
				continue
			}

			err = f.SetModel(model)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fileName, err := f.WriteModel(outputDir, framework.InitialBuffers{})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("\t- written to %s\n", fileName)
		}
	}

	return
}

func runCode(frameworks framework.List) {
	for _, f := range frameworks {
		result, err := f.Run(framework.InitialBuffers{})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("== %s ==\n", f.Info().Name)
		fmt.Println(string(result.Output))
		fmt.Println()
	}
}
