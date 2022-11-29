package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/util/chalk"
)

var (
	ErrRequiresSubcommand = errors.New("env command requires subcommand")
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Setup & maintain an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = ErrRequiresSubcommand
		chalk.PrintErr(err)

		return
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
