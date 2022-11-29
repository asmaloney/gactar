package cmd

import (
	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/amod"
)

var ebnfCmd = &cobra.Command{
	Use:   "ebnf",
	Short: "Output amod EBNF to stdout and quit",
	Run: func(cmd *cobra.Command, args []string) {
		amod.OutputEBNF()
	},
}

func init() {
	rootCmd.AddCommand(ebnfCmd)
}
