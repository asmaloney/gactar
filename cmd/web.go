package cmd

import (
	"github.com/spf13/cobra"

	"github.com/asmaloney/gactar/examples"
	"github.com/asmaloney/gactar/modes/web"
)

var (
	flagPort = 8181
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start a web server to run in a browser",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		settings, err := setupForRun(cmd)
		if err != nil {
			return err
		}

		w, err := web.Initialize(settings, flagPort, &examples.AMODExamples)
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
	rootCmd.AddCommand(webCmd)

	webCmd.Flags().IntVarP(&flagPort, "port", "p", flagPort, "port to run the web server on")
}
