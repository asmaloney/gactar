// Package clicontext provide support functions for working with cli.Context.
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/filesystem"
)

var (
	ErrInvalidContext = errors.New("internal error: invalid context")
)

type Settings struct {
	EnvPath  string // full path to the environment directory
	TempPath string // full path to the temp directory

	ActiveFrameworks framework.List // active frameworks (set from the command line)

	Version string // the version string for output to command line
}

// SetupPaths will set PATH and VIRTUAL_ENV environment variables to our environment path.
func SetupPaths(envPath string) (err error) {
	envPath, err = filepath.Abs(envPath)
	if err != nil {
		return
	}

	// Skip if we have already set up with this env
	currentVirtEnv := os.Getenv("VIRTUAL_ENV")
	if currentVirtEnv == envPath {
		return nil
	}

	pythonVENVPath := "bin"

	// Python on Windows puts itself in a different place, so adjust our path accordingly
	if runtime.GOOS == "windows" {
		pythonVENVPath = "Scripts"
	}

	pythonVENVPath = filepath.Join(envPath, pythonVENVPath)
	cclPath := filepath.Join(envPath, "ccl")

	// Restrict PATH to our envPath's "bin" dir & the "ccl" dir
	err = os.Setenv("PATH", fmt.Sprintf("%s%c%s", cclPath, os.PathListSeparator, pythonVENVPath))
	if err != nil {
		return
	}

	err = os.Setenv("VIRTUAL_ENV", envPath)
	if err != nil {
		return
	}

	os.Unsetenv("PYTHONHOME")

	return
}

func CreateTempDir(settings *Settings) (path string, err error) {
	path = settings.TempPath
	if path == "" {
		path = fmt.Sprintf("%s%c%s", os.Getenv("VIRTUAL_ENV"), filepath.Separator, "gactar-temp")
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return
	}

	err = filesystem.CreateDir(path)
	if err != nil {
		return
	}

	settings.TempPath = path
	return
}
