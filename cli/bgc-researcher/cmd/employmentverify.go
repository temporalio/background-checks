package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/utils"
	"github.com/temporalio/background-checks/workflows"
)

var employmentVerifyCmd = &cobra.Command{
	Use:   "employmentverify",
	Short: "Complete the employment verification process for a candidate",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		router := api.Router(nil)

		requestURL, err := router.Get("employmentverify").Host(APIEndpoint).URL("token", Token)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		submission := workflows.EmploymentVerificationSubmissionSignal{
			EmploymentVerificationComplete: true,
			EmployerVerified:               true,
		}

		response, err := utils.PostJSON(requestURL, submission)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		if response.StatusCode != http.StatusOK {
			log.Fatalf("%s: %s", http.StatusText(response.StatusCode), body)
		}
		fmt.Println("Employment verification received")
	},
}

func init() {
	rootCmd.AddCommand(employmentVerifyCmd)
	employmentVerifyCmd.Flags().StringVar(&Token, "token", "", "Token")
	employmentVerifyCmd.MarkFlagRequired("token")

}
