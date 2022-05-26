package filesystem

import (
	"os"

	"github.com/urfave/cli/v2"
)

func CreateDir(path string) (err error) {
	err = os.MkdirAll(path, 0750)
	if err != nil && !os.IsExist(err) {
		return
	}

	return
}

func CreateTempDir(ctx *cli.Context) (err error) {
	path := ctx.Path("temp")
	return CreateDir(path)
}
