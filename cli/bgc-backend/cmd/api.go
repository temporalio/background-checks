package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"

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
			log.Fatalf("error: %v", err)
		}
		defer c.Close()

		srv := &http.Server{
			Handler: api.Router(c),
			Addr:    api.DefaultEndpoint,
		}

		errCh := make(chan error, 1)
		go func() { errCh <- srv.ListenAndServe() }()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case <-sigCh:
			srv.Close()
		case err = <-errCh:
			log.Fatalf("error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
