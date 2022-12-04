package cmd

import (
	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/modes/shell"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run an interactive shell",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		settings, err := setupForRun(cmd)
		if err != nil {
			return err
		}

		w, err := shell.Initialize(settings)
		if err != nil {
			return err
		}

		err = w.Start()
		if err != nil {
			return err
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)
}
