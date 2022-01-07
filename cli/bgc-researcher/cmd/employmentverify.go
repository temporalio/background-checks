/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

		submission := types.EmploymentVerificationSubmissionSignal{
			EmploymentVerificationComplete: true,
			EmployerVerified:               true,
		}

		response, err := utils.PostJSON(requestURL, submission)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)

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
