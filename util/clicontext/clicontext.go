// Package clicontext provide support functions for working with cli.Context.
package clicontext

import (
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
	path, err := ExpandPath(ctx, "temp")
	if err != nil {
		return
	}

	return filesystem.CreateDir(path)
}
