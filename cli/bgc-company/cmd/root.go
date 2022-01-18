package cmd

import (
	"github.com/spf13/cobra"
)

const APIEndpoint = "api:8081"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bgc-company",
	Short: "Company CLI for the Background Check application",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
