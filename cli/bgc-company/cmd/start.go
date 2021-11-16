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

var (
	tier string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [email]",
	Short: "starts a background check for a candidate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

		router := api.Router()

		requestURL, err := router.Get("checks_create").Host(api.DefaultEndpoint).URL()
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		input := types.BackgroundCheckWorkflowInput{
			Email: email,
			Tier:  tier,
		}

		response, err := utils.PostJSON(requestURL, input)
		if err != nil {
			log.Fatalf("request error: %v", err)
		}

		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != http.StatusCreated {
			log.Fatalf("%s: %s", http.StatusText(response.StatusCode), body)
		}

		fmt.Printf("Created check")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&tier, "tier", "standard", "Tier (which checks to run)")
}
