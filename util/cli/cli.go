// Package clicontext provide support functions for working with cli.Context.
package cli

import (
	"context"
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

	Frameworks framework.List // active frameworks

	Version string // the version string for output to command line
}
type settingsKey string

var contextKey settingsKey = "cli"

func NewContext(ctx context.Context, u *Settings) context.Context {
	return context.WithValue(ctx, contextKey, u)
}

func FromContext(ctx context.Context) (*Settings, error) {
	u, ok := ctx.Value(contextKey).(*Settings)

	if !ok {
		return nil, ErrInvalidContext
	}

	return u, nil
}

// SetupPaths will set PATH and VIRTUAL_ENV environment variables to our environment path.
func SetupPaths(envPath string) (err error) {
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
