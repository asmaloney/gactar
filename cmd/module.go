package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/asmaloney/gactar/actr/modules"
	"github.com/asmaloney/gactar/actr/param"
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
		fmt.Println(chalk.BoldHeader(" Description"))
		fmt.Printf("   %s\n", mod.ModuleDescription())

		writer := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

		if mod.Parameters() != nil {
			fmt.Fprintln(writer, chalk.BoldHeader(" Module Parameters"))

			params := mod.Parameters().ParameterList()
			outputParams(writer, 1, params)
		}

		if mod.HasBuffers() {
			fmt.Fprintln(writer, chalk.BoldHeader(" Buffers"))

			for _, buffer := range mod.Buffers() {
				fmt.Fprintf(writer, "   %s\n", chalk.Italic(buffer.Name()))

				if buffer.RequestParameters() != nil {
					fmt.Fprintf(writer, "\t\t%s\n", chalk.Header("Request Parameters:"))

					params := buffer.RequestParameters().ParameterList()
					outputParams(writer, 3, params)
				}

				if buffer.Parameters() != nil {
					fmt.Fprintf(writer, "\t\t%s\n", chalk.Header("Config Options:"))

					params := buffer.Parameters().ParameterList()
					outputParams(writer, 3, params)

				}
			}
		}

		fmt.Fprintln(writer, "")
		writer.Flush()
	}
}

func outputParams(writer *tabwriter.Writer, level int, list param.List) {
	for _, param := range list {
		outputParam(writer, level, param)
	}
}

func outputParam(writer *tabwriter.Writer, level int, p param.ParamInterface) {
	tabs := strings.Repeat("\t", level)

	fmt.Fprintf(writer, "%s%s", tabs, chalk.Italic(p.Name()))

	var typeStr string
	var valuesStr string

	switch v := p.(type) {
	case param.Bool:
		{
			typeStr = "bool"
			valuesStr = fmt.Sprintf("%strue,false", tabs)
		}

	case param.Str:
		{
			typeStr = "string"
			valuesStr = fmt.Sprintf("%s%s", tabs, strings.Join(v.ValidValues(), ","))
		}

	case param.Int:
		{
			var minStr string
			var maxStr string

			typeStr = "int"
			if v.Min() != nil {
				minStr = strconv.Itoa(*v.Min())
			}
			if v.Max() != nil {
				maxStr = strconv.Itoa(*v.Max())
			}

			valuesStr = fmt.Sprintf("%s%s-%s", tabs, minStr, maxStr)
		}

	case param.Float:
		{
			var minStr string
			var maxStr string

			typeStr = "float"
			if v.Min() != nil {
				minStr = strconv.FormatFloat(*v.Min(), 'f', 2, 64)
			}
			if v.Max() != nil {
				maxStr = strconv.FormatFloat(*v.Max(), 'f', 2, 64)
			}

			valuesStr = fmt.Sprintf("%s%s-%s", tabs, minStr, maxStr)
		}
	}

	fmt.Fprintf(writer, "%s%s", tabs, typeStr)
	fmt.Fprint(writer, valuesStr)
	fmt.Fprintf(writer, "%s%s", tabs, p.Description())

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
