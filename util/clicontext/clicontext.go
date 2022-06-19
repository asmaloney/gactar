// Package clicontext provide support functions for working with cli.Context.
package clicontext

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/urfave/cli/v2"
)

// SetupPaths will set PATH and VIRTUAL_ENV environment variables to our environment path.
func SetupPaths(envPath string) (err error) {
	pythonVENVPath := "bin"

	// Python on Windows puts itself in a different place, so adjust our path accordingly
	if runtime.GOOS == "windows" {
		pythonVENVPath = "Script"
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

// ExpandPath expands the given path and sets it back in the context.
func ExpandPath(ctx *cli.Context, flag string) (path string, err error) {
	path, err = filepath.Abs(ctx.Path(flag))
	if err != nil {
		return
	}

	err = ctx.Set(flag, path)
	if err != nil {
		return "", err
	}

	return
}

// CreateTempDir looks up the "temp" flag in our context, expands the path, and creates the dir.
func CreateTempDir(ctx *cli.Context) (err error) {
	if !ctx.IsSet("temp") {
		defaultTemp := fmt.Sprintf("%s%c%s", os.Getenv("VIRTUAL_ENV"), filepath.Separator, "gactar-temp")
		err = ctx.Set("temp", defaultTemp)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
	}

	path, err := ExpandPath(ctx, "temp")
	if err != nil {
		return
	}

	return filesystem.CreateDir(path)
}
