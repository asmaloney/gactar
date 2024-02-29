// Package frameworkutil provides functions to work with frameworks. It's a separate utility
// to avoid circular dependencies with the framework implementations.
package frameworkutil

import (
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/framework/ccm_pyactr"
	"github.com/asmaloney/gactar/framework/pyactr"
	"github.com/asmaloney/gactar/framework/vanilla_actr"

	"github.com/asmaloney/gactar/util/chalk"
	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/runoptions"
)

// CreateFrameworks takes a slice of framework names and some settings,
// creates any valid ones, and returns a list of them.
// If "names" is empty it will try to create all valid frameworks.
func CreateFrameworks(settings *cli.Settings, names []string) (list framework.List) {
	if len(names) == 0 {
		names = runoptions.ValidNamedFrameworks()
	}

	list = framework.List{}

	for _, f := range names {
		var fw framework.Framework
		var err error

		switch f {
		case "ccm":
			fw, err = ccm_pyactr.New(settings.TempPath)

		case "pyactr":
			fw, err = pyactr.New(settings.TempPath)

		case "vanilla":
			fw, err = vanilla_actr.New(settings.TempPath)

		default:
			chalk.PrintErrStr("unknown framework:", f)
			continue
		}

		if err != nil {
			chalk.PrintErr(err)
			continue
		}

		list[f] = fw
	}

	return
}
