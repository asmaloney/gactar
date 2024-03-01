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
	"github.com/asmaloney/gactar/util/validate"
)

var (
	ErrNoInputFiles     = errors.New("no input files specified on command line")
	ErrNoFilesToProcess = errors.New("no files to process")
	ErrNoValidModels    = errors.New("no valid models to run")
)

type DefaultMode struct {
	settings *cli.Settings

	runAfterGenerate bool
	fileList         []string
}

func Initialize(settings *cli.Settings, files []string, runAfterGenerate bool) (d *DefaultMode, err error) {
	d = &DefaultMode{
		settings:         settings,
		runAfterGenerate: runAfterGenerate,
	}

	// Check if files exist first
	if len(files) == 0 {
		return nil, ErrNoInputFiles
	}

	existingFiles := files[:0]
	for _, file := range files {
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

	d.fileList = existingFiles
	return
}

func (d *DefaultMode) Start() (err error) {
	fmt.Printf("Intermediate file path: %q\n", d.settings.TempPath)

	err = generateCode(d.settings.Frameworks, d.fileList, d.settings.TempPath)
	if err != nil {
		return err
	}

	if d.runAfterGenerate {
		runCode(d.settings.Frameworks)
	}
	return
}

func generateCode(frameworks framework.List, files []string, outputDir string) (err error) {
	modelMap := map[string]*actr.Model{}

	for _, file := range files {
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

	for _, f := range frameworks {
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

			fileName, err := f.WriteModel(outputDir, &model.DefaultParams, framework.InitialBuffers{})
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
		model := f.Model()
		result, err := f.Run(&model.DefaultParams, framework.InitialBuffers{})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("== %s ==\n", f.Info().Name)
		fmt.Println(string(result.Output))
		fmt.Println()
	}
}
