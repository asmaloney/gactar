// Package clicontext provide support functions for working with cli.Context.
package clicontext

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/urfave/cli/v2"
)

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
