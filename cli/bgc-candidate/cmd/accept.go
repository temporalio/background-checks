package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/cli/utils"
	"github.com/temporalio/background-checks/types"
)

// acceptCmd represents the accept command
var acceptCmd = &cobra.Command{
	Use:   "accept",
	Short: "Accept a background check",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		router := api.Router(nil)

		requestURL, err := router.Get("accept").Host(APIEndpoint).URL("token", Token)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		candidatedetails := types.CandidateDetails{
			FullName: FullName,
			SSN:      SSN,
			Employer: Employer}
		submission := types.AcceptSubmissionSignal{
			CandidateDetails: candidatedetails,
		}

		response, err := utils.PostJSON(requestURL, submission)
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != http.StatusOK {
			log.Fatalf("%s: %s", http.StatusText(response.StatusCode), body)
		}

		fmt.Println("Accepted")
	},
}

func init() {
	rootCmd.AddCommand(acceptCmd)
	acceptCmd.Flags().StringVar(&Token, "token", "", "Token")
	acceptCmd.MarkFlagRequired("token")
	acceptCmd.Flags().StringVar(&FullName, "fullname", "", "Candidate's full name")
	acceptCmd.MarkFlagRequired("fullname")
	acceptCmd.Flags().StringVar(&SSN, "ssn", "", "Social Security #")
	acceptCmd.MarkFlagRequired("ssn")
	acceptCmd.Flags().StringVar(&Employer, "employer", "", "Social Security #")
	acceptCmd.MarkFlagRequired("employer")

}
