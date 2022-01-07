package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/cli/utils"
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "cancels a background check",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		router := api.Router(nil)

		requestURL, err := router.Get("check_cancel").Host(APIEndpoint).URL("email", email, "id", id)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		response, err := utils.PostJSON(requestURL, nil)
		if err != nil {
			log.Fatalf("request error: %v", err)
		}
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != http.StatusOK {
			log.Fatalf("%s: %s", http.StatusText(response.StatusCode), body)
		}

		fmt.Printf("Cancelled check\n")
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().StringVar(&email, "email", "", "Candidate's email address")
	cancelCmd.MarkFlagRequired("email")
	cancelCmd.Flags().StringVar(&id, "id", "", "Check ID")
	cancelCmd.MarkFlagRequired("id")
}
