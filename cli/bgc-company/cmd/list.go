package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/cli/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List background checks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		router := api.Router(nil)

		requestURL, err := router.Get("checks_create").Host(APIEndpoint).URL()
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		query := requestURL.Query()
		query.Set("email", email)
		query.Set("status", status)
		requestURL.RawQuery = query.Encode()

		var checks []api.BackgroundCheck
		_, err = utils.GetJSON(requestURL, &checks)
		if err != nil {
			log.Fatalf("request error: %v", err)
		}

		fmt.Printf("Background Checks:\n")
		for _, check := range checks {
			fmt.Printf("ID: %s Email: %s Status: %s\n", check.ID, check.Email, check.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&email, "email", "", "Candidate's email address")
	listCmd.Flags().StringVar(&status, "status", "", "Status")
}
