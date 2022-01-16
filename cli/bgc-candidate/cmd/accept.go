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
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/utils"
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

		body, _ := io.ReadAll(response.Body)

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
