// Package defaultmode is used for running gactar models on the command line.
package defaultmode

import (
	"errors"
	"fmt"
	"os"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/runoptions"
	"github.com/asmaloney/gactar/util/validate"
)

var (
	ErrNoInputFiles     = errors.New("no input files specified on command line")
	ErrNoFilesToProcess = errors.New("no files to process")
	ErrNoValidModels    = errors.New("no valid models to run")
)

// CommandLineOptions come from the command line.
type CommandLineOptions struct {
	FileList           []string
	RunAfterGeneration bool

	// these override any options from the model
	runoptions.Options
}

type DefaultMode struct {
	settings *cli.Settings

	commandLineOptions CommandLineOptions
}

func Initialize(settings *cli.Settings, options CommandLineOptions) (d *DefaultMode, err error) {
	// Check if files exist first
	if len(options.FileList) == 0 {
		return nil, ErrNoInputFiles
	}

	var existingFiles []string

	for _, file := range options.FileList {
		if _, fileErr := os.Stat(file); errors.Is(fileErr, os.ErrNotExist) {
			fileErr = &filesystem.ErrFileDoesNotExist{FileName: file}
			chalk.PrintErr(fileErr)

			continue
		}

		existingFiles = append(existingFiles, file)
	}

	if len(existingFiles) == 0 {
		return nil, ErrNoFilesToProcess
	}

	d = &DefaultMode{
		settings:           settings,
		commandLineOptions: options,
	}

	d.commandLineOptions.FileList = existingFiles

	return
}

func (d *DefaultMode) Start() (err error) {
	fmt.Printf("Intermediate file path: %q\n", d.settings.TempPath)

	err = d.generateCode()
	if err != nil {
		return err
	}

	if d.commandLineOptions.RunAfterGeneration {
		d.runCode(d.settings.ActiveFrameworks)
	}
	return
}

func (d *DefaultMode) generateCode() (err error) {
	modelMap := map[string]*actr.Model{}

	for _, file := range d.commandLineOptions.FileList {
		fmt.Printf("Generating model for %s\n", file)
		model, log, modelErr := amod.GenerateModelFromFile(file)
		if modelErr != nil {
			fmt.Print(log)
			continue
		}

		// When using "-r" the goal must be initialized in the code.
		validate.Goal(model, "", log)

		fmt.Print(log)

		modelMap[file] = model
	}

	if len(modelMap) == 0 {
		return ErrNoValidModels
	}

	for _, f := range d.settings.ActiveFrameworks {
		fmt.Printf(" %s\n", f.Info().Name)
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

			options := overrideRunOptions(&model.DefaultParams, &d.commandLineOptions.Options)

			fileName, err := f.WriteModel(d.settings.TempPath, options)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("\t- written to %s\n", fileName)
		}
	}

	return
}

func (d *DefaultMode) runCode(frameworks framework.List) {
	for _, f := range frameworks {
		model := f.Model()

		options := overrideRunOptions(&model.DefaultParams, &d.commandLineOptions.Options)

		result, err := f.Run(options)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("== %s ==\n", f.Info().Name)
		fmt.Println(string(result.Output))
		fmt.Println()
	}
}

// overrideRunOptions overrides options set in the model with any set on the command line.
func overrideRunOptions(modelOptions, cliOptions *runoptions.Options) *runoptions.Options {
	options := *modelOptions

	if cliOptions.LogLevel != nil {
		options.LogLevel = cliOptions.LogLevel
	}

	if cliOptions.TraceActivations != nil {
		options.TraceActivations = cliOptions.TraceActivations
	}

	if cliOptions.RandomSeed != nil {
		options.RandomSeed = cliOptions.RandomSeed
	}

	return &options
}
