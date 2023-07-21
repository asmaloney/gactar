package cmd

import (
	"fmt"
	"os"
	"strconv"
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
			params := mod.Parameters()

			writer := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
			for _, param := range params {
				outputParam(writer, param)
			}
			writer.Flush()
		}

		fmt.Println("")
	}
}

func outputParam(writer *tabwriter.Writer, param modules.ParamInterface) {
	fmt.Fprintf(writer, "\t%s", chalk.Italic(param.GetName()))

	var typeStr string
	var minStr string
	var maxStr string

	switch v := param.(type) {
	case modules.ParamInt:
		{
			typeStr = "int"
			if v.Min != nil {
				minStr = strconv.Itoa(*v.Min)
			}
			if v.Max != nil {
				maxStr = strconv.Itoa(*v.Max)
			}
		}
	case modules.ParamFloat:
		{
			typeStr = "float"
			if v.Min != nil {
				minStr = strconv.FormatFloat(*v.Min, 'f', 2, 64)
			}
			if v.Max != nil {
				maxStr = strconv.FormatFloat(*v.Max, 'f', 2, 64)
			}
		}
	}

	fmt.Fprintf(writer, "\t%s", typeStr)
	fmt.Fprintf(writer, "\t%s-%s", minStr, maxStr)
	fmt.Fprintf(writer, "\t%s", param.GetDescription())

	fmt.Fprintln(writer, "")
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
