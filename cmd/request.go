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
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/client"
)

var (
	email string
	tier  string
)

// requestCmd represents the start command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Requests a background check for a candidate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient(client.Options{})
		if err != nil {
			log.Fatalf("client error: %v", err)
		}
		defer c.Close()

		_, err = c.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				TaskQueue: "background-checks-main",
			},
			workflows.BackgroundCheck,
			internal.BackgroundCheckInput{
				Email: email,
				Tier:  tier,
			},
		)

		if err != nil {
			log.Fatalf("Failed to start background check: %v", err)
		}

		log.Print("Background check requested")
	},
}

func init() {
	uiCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringVar(&email, "email", "", "Email address for candidate")
	requestCmd.Flags().StringVar(&tier, "tier", "", "Check tier")
}
