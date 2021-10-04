package main

import (
	"embed"
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
		HideHelpCommand:      true,
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

			&cli.BoolFlag{Name: "interactive", Aliases: []string{"i"}, Usage: "run an interactive shell"},

			&cli.BoolFlag{Name: "web", Aliases: []string{"w"}, Usage: "start a webserver to run in a browser"},
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: defaultPort, Usage: "port to run the webserver on"},
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

			if c.Bool("web") {
				w, err := web.Initialize(c, frameworks, &amodExamples)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				err = w.Start()
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
			}

			if c.Int("port") != defaultPort {
				fmt.Println("info: -port only applies when using -web")
			}

			if c.Bool("interactive") {
				s, err := shell.Initialize(c, &frameworks)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				err = s.Start()
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				return nil
			}

			// We are not interactive or web, so simply generate the output files.

			cli.ShowVersion(c)

			generateCode(&frameworks, c.Args().Slice())
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
			fmt.Println(err.Error())
		}
	}

	if len(frameworks) == 0 {
		err = fmt.Errorf("could not create any frameworks - please check your installation")
		return framework.List{}, err
	}

	return
}

func generateCode(frameworks *framework.List, files []string) {
	var err error

	if len(files) == 0 {
		err = fmt.Errorf("no input files specified on command line")
		fmt.Println(err.Error())
		return
	}

	for _, framework := range *frameworks {
		err = framework.Initialize()
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

			err = framework.SetModel(model)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fileName, err := framework.WriteModel("", "")
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("\t- written to %s\n", fileName)
		}
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
