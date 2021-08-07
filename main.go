package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"gitlab.com/asmaloney/gactar/amod"
	"gitlab.com/asmaloney/gactar/framework/pyactr"
	"gitlab.com/asmaloney/gactar/shell"
	"gitlab.com/asmaloney/gactar/web"
)

func main() {
	defaultPort := 8181

	app := &cli.App{
		Name:        "gactar",
		Usage:       "A command-line tool for working with ACT-R models",
		Description: "A proof-of-concept tool for creating ACT-R models using a declarative file format.",
		Version:     "v0.0.2",
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
			&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Usage: "turn on debugging output"},
			&cli.BoolFlag{Name: "ebnf", Usage: "output amod EBNF to stdout and quit"},

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

			if c.Bool("web") {
				w, err := web.Initialize(c)
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
				s, err := shell.Initialize(c)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				err = s.Start()
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
			}

			cli.ShowVersion(c)

			framework, err := pyactr.Initialize()
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			for _, arg := range c.Args().Slice() {
				fmt.Printf("-- Generating code for %s\n", arg)
				model, err := amod.GenerateModelFromFile(arg)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				err = framework.SetModel(model)
				if err != nil {
					fmt.Println(err.Error())
					return err
				}

				fileName, err := framework.WriteModel("")
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
				fmt.Printf("   Written to %s\n", fileName)
			}

			return nil
		},
	}

	// Used to output command line options for documentation.
	// fmt.Println(app.ToMarkdown())

	app.Run(os.Args)
}
