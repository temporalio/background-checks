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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [email]",
	Short: "List background checks which need consent.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

		router := api.Router()

		requestURL, err := router.Get("consents").Host(api.DefaultEndpoint).URL("email", email)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		var consents []types.CandidateTodo
		response, err := utils.GetJSON(requestURL, &consents)

		if response.StatusCode == http.StatusNotFound {
			fmt.Printf("No outstanding consents for: %s\n", email)
			return
		}
		if err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Printf("Consents:\n")
		for i, consent := range consents {
			fmt.Printf("%d: token: %s\n", i, consent.Token)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
