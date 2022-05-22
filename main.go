package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/framework/ccm_pyactr"
	"github.com/asmaloney/gactar/framework/pyactr"
	"github.com/asmaloney/gactar/framework/vanilla_actr"
	"github.com/asmaloney/gactar/shell"
	"github.com/asmaloney/gactar/version"
	"github.com/asmaloney/gactar/web"
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

			&cli.StringSliceFlag{
				Name:    "framework",
				Aliases: []string{"f"},
				Value:   defaultFramework,
				Usage:   fmt.Sprintf("add framework - valid frameworks: %s", strings.Join(framework.ValidFrameworks, ", ")),
			},

			// for default
			&cli.PathFlag{Name: "output", Aliases: []string{"o"}, Value: ".", Usage: "directory for generated files (will be created)"},
			&cli.BoolFlag{Name: "run", Aliases: []string{"r"}, Usage: "run the models after generating the code"},

			// for interactive
			&cli.BoolFlag{Name: "interactive", Aliases: []string{"i"}, Usage: "run an interactive shell"},

			// for web
			&cli.BoolFlag{Name: "web", Aliases: []string{"w"}, Usage: "start a web server to run in a browser"},
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: defaultPort, Usage: "port to run the web server on"},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				amod.SetDebug(true)
			}

			if c.Bool("ebnf") {
				amod.OutputEBNF()
				return nil
			}

			frameworks, err := createFrameworks(c)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			if c.Bool("web") && c.Bool("interactive") {
				err = errors.New("cannot run 'web' and 'interactive' at the same time")
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

	list = deduplicate(list)
	sort.Strings(list)

	if list[0] == "all" {
		list = framework.ValidFrameworks[1:]
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

func handleWeb(context *cli.Context, frameworks framework.List) (err error) {
	w, err := web.Initialize(context, frameworks, &amodExamples)
	if err != nil {
		return err
	}

	err = w.Start()
	if err != nil {
		return err
	}

	return
}

func handleInteractive(context *cli.Context, frameworks framework.List) (err error) {
	s, err := shell.Initialize(context, frameworks)
	if err != nil {
		return err
	}

	err = s.Start()
	if err != nil {
		return err
	}

	return
}

func handleDefault(context *cli.Context, frameworks framework.List) (err error) {
	cli.ShowVersion(context)

	// Check if files exist first
	files := context.Args().Slice()

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

	outputDir := context.Path("output")
	err = os.MkdirAll(outputDir, 0750)
	if err != nil && !os.IsExist(err) {
		return
	}

	generateCode(frameworks, existingFiles, outputDir, context.Bool("run"))
	if err != nil {
		return err
	}

	if context.Bool("run") {
		runCode(frameworks)
	}
	return
}

func generateCode(frameworks framework.List, files []string, outputDir string, runCode bool) {
	var err error

	for _, f := range frameworks {
		err = f.Initialize()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for _, file := range files {
			fmt.Printf("\t- Generating code for %s\n", file)
			model, log, err := amod.GenerateModelFromFile(file)
			fmt.Print(log)
			if err != nil {
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
}

func runCode(frameworks framework.List) {
	for _, f := range frameworks {
		_, result, err := f.Run(framework.InitialBuffers{})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("== %s ==\n", f.Info().Name)
		fmt.Println(string(result))
		fmt.Println()
	}
}

func deduplicate(s []string) []string {
	if len(s) <= 1 {
		return s
	}

	result := []string{}
	seen := make(map[string]struct{})
	for _, val := range s {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = struct{}{}
		}
	}

	return result
}
