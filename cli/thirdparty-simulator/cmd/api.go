package cmd

import (
	"github.com/spf13/cobra"

	"github.com/temporalio/background-checks/cli/thirdparty-simulator/api"
)

// apiCmd represents the thirdparty command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts an API server for the third party API simulator",
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
