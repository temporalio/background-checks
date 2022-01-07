package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/temporal"
	"go.temporal.io/sdk/client"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run API Server",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := temporal.NewClient(client.Options{})
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		srv := &http.Server{
			Handler: api.Router(c),
			Addr:    "0.0.0.0:8081",
		}

		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
