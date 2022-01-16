/*
 * The MIT License
 *
 * Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
 *
 * Copyright (c) 2020 Uber Technologies, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List background checks",
	Args:  cobra.NoArgs,
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
