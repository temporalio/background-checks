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
	"github.com/temporalio/background-checks/mocks"
)

// consentCmd represents the consent command
var consentCmd = &cobra.Command{
	Use:   "consent [token]",
	Short: "Consent to a background check",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]

		router := api.Router()

		requestURL, err := router.Get("consent").Host(api.DefaultEndpoint).URL("token", token)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		response, err := utils.PostJSON(requestURL, mocks.ConsentResultConsented)
		if err != nil {
			log.Fatalf(err.Error())
		}

		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != http.StatusOK {
			log.Fatalf("%s: %s", http.StatusText(response.StatusCode), body)
		}

		fmt.Println("Recorded consent")
	},
}

func init() {
	rootCmd.AddCommand(consentCmd)
}
