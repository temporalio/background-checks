package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/ui"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Run UI Server",
	Run: func(cmd *cobra.Command, args []string) {
		srv := &http.Server{
			Handler: ui.Router(),
			Addr:    "0.0.0.0:8083",
		}

		errCh := make(chan error, 1)
		go func() { errCh <- srv.ListenAndServe() }()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case <-sigCh:
			srv.Close()
		case err := <-errCh:
			log.Fatalf("error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
