package framework

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"

	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/python"
)

// Some tools for working with our frameworks

// Reads a file, generates a model, validates it, and generates code from it for a given framework.
// This is useful for testing.
func GenerateCodeFromFile(fw Framework, inputFile string, initialBuffers InitialBuffers) (code []byte, err error) {
	amodCode, err := os.ReadFile(inputFile)
	if err != nil {
		return
	}

	model, log, err := amod.GenerateModel(string(amodCode))
	if err != nil {
		fmt.Print(log)
		return
	}

	log = fw.ValidateModel(model)
	if log.HasIssues() {
		if log.HasError() {
			err = &ErrModelValidationFailed{Log: log}
			return
		}
	}

	err = fw.SetModel(model)
	if err != nil {
		return
	}

	code, err = fw.GenerateCode(initialBuffers)
	if err != nil {
		return
	}

	return
}

// Setup will check that the executable exists and then use it to identify itself.
func Setup(info *Info) (err error) {
	_, err = filesystem.CheckForExecutable(info.ExecutableName)
	if err != nil {
		return
	}

	err = identifyYourself(info.Name, info.ExecutableName)
	if err != nil {
		return
	}

	for _, packageName := range info.PythonRequiredPackages {
		err = python.CheckForPackage(info.ExecutableName, packageName)
		if err != nil {
			return
		}
	}

	return
}

func ParseInitialBuffers(model *actr.Model, initialBuffers InitialBuffers) (parsed ParsedInitialBuffers, err error) {
	parsed = ParsedInitialBuffers{}

	for bufferName, bufferInit := range initialBuffers {
		buffer := model.LookupBuffer(bufferName)
		if buffer == nil {
			err = &ErrBufferNotFound{
				BufferName: bufferName,
				ModelName:  model.Name,
			}
			return
		}

		pattern, parseErr := amod.ParseChunk(model, bufferInit)
		if parseErr != nil {
			err = fmt.Errorf("in initial buffer '%s' - %w", bufferName, parseErr)
			return
		}

		parsed[bufferName] = pattern
	}

	return
}

// identifyYourself outputs version info for an executable.
func identifyYourself(frameworkName, exeName string) (err error) {
	cmd := exec.Command(exeName, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	version := strings.TrimSpace(string(output))

	fmt.Printf("%s: Using %s\n", frameworkName, version)

	return
}
