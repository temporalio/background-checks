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
	"net/http"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/cli/utils"
	"github.com/temporalio/background-checks/types"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [email]",
	Short: "show background check status.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

		router := api.Router()

		requestURL, err := router.Get("candidate").Host(api.DefaultEndpoint).URL("email", email)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		var check types.BackgroundCheckStatusSignal
		response, err := utils.GetJSON(requestURL, &check)

		if response.StatusCode == http.StatusNotFound {
			fmt.Printf("No background check found for: %s\n", email)
			return
		}
		if err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Printf("Status: %s\n", check.Status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
