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
	"log"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/cli/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List background checks",
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
