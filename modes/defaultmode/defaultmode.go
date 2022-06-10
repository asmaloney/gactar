package defaultmode

import (
	"errors"
	"fmt"
	"os"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/util/validate"
	"github.com/urfave/cli/v2"
)

type DefaultMode struct {
	context        *cli.Context
	actrFrameworks framework.List

	fileList []string

	tempPath string
}

func Initialize(ctx *cli.Context, frameworks framework.List) (d *DefaultMode, err error) {
	d = &DefaultMode{
		context:        ctx,
		actrFrameworks: frameworks,

		tempPath: ctx.Path("temp"),
	}

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

	d.fileList = existingFiles
	return
}

func (d *DefaultMode) Start() (err error) {
	fmt.Printf("Intermediate file path: %q\n", d.tempPath)

	err = generateCode(d.actrFrameworks, d.fileList, d.tempPath, d.context.Bool("run"))
	if err != nil {
		return err
	}

	if d.context.Bool("run") {
		runCode(d.actrFrameworks)
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
