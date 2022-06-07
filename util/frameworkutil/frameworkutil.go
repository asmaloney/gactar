// Package frameworkutil provides functions to work with frameworks. It's a separate utility
// to avoid circular dependencies with the framework implementations.
package frameworkutil

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/framework/ccm_pyactr"
	"github.com/asmaloney/gactar/framework/pyactr"
	"github.com/asmaloney/gactar/framework/vanilla_actr"
)

// CreateFrameworks takes a slice of framework names and a context and will
// create any valid ones and return a list of them. If "names" is empty it will
// try to create all valid frameworks.
func CreateFrameworks(ctx *cli.Context, names []string) (list framework.List) {
	if len(names) == 0 {
		names = framework.ValidNamedFrameworks()
	}

	list = framework.List{}

	for _, f := range names {
		var createErr error
		switch f {
		case "ccm":
			list["ccm"], createErr = ccm_pyactr.New(ctx)
		case "pyactr":
			list["pyactr"], createErr = pyactr.New(ctx)
		case "vanilla":
			list["vanilla"], createErr = vanilla_actr.New(ctx)
		default:
			fmt.Printf("unknown framework: %s\n", f)
		}

		if createErr != nil {
			fmt.Println(createErr.Error())
		}
	}

	return
}
