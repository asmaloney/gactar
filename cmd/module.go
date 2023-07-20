package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/asmaloney/gactar/actr/modules"
	"github.com/asmaloney/gactar/util/chalk"
	"github.com/spf13/cobra"
)

var modulesCmd = &cobra.Command{
	Use:   "module",
	Short: "Get info about available modules",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = errRequiresSubcommand{command: cmd}
		return
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get detailed info about modules by name (list of names or 'all')",
	Run: func(cmd *cobra.Command, args []string) {
		info(args)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output module names and their descriptions",
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

type errUnknownModule struct {
	module string
}

func (e errUnknownModule) Error() string {
	message := fmt.Sprintf("module not found: %q", e.module)
	return chalk.ErrorBold(message)
}

func init() {
	modulesCmd.AddCommand(infoCmd)
	modulesCmd.AddCommand(listCmd)

	rootCmd.AddCommand(modulesCmd)
}

func info(args []string) {
	if len(args) == 1 && args[0] == "all" {
		args = modules.ModuleNames()
	}

	for _, name := range args {
		mod := modules.FindModule(name)
		if mod == nil {
			err := errUnknownModule{module: name}
			chalk.PrintErr(err)

			continue
		}

		fmt.Printf("%s: %s (%s)\n", chalk.Bold("Module"), mod.ModuleName(), mod.ModuleVersion())
		fmt.Println(chalk.Header(" Description"))
		fmt.Printf("  %s\n", mod.ModuleDescription())

		if mod.HasParameters() {
			fmt.Println(chalk.Header(" Parameters"))
			writer := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)

			for _, paramName := range mod.ParameterNames() {
				fmt.Fprintf(writer, "\t%s", chalk.Italic(paramName))

				info := mod.ParameterInfo(paramName)
				if info != nil {
					fmt.Fprintf(writer, "\t%s", info.Description)
				}

				fmt.Fprintln(writer, "")

			}
			writer.Flush()
		}

		fmt.Println("")
	}
}

func list() {
	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)

	for _, module := range modules.AllModules() {
		version := module.ModuleVersion()

		if version == modules.BuiltIn {
			version = chalk.Italic(version)
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\n", chalk.Bold(module.ModuleName()), version, module.ModuleDescription())
	}

	writer.Flush()
}
