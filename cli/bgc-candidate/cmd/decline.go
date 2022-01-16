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
	"github.com/temporalio/background-checks/utils"
)

// declineCmd represents the decline command
var declineCmd = &cobra.Command{
	Use:   "decline",
	Short: "Decline a background check",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		router := api.Router(nil)

		requestURL, err := router.Get("decline").Host(APIEndpoint).URL("token", Token)
		if err != nil {
			log.Fatalf("cannot create URL: %v", err)
		}

		response, err := utils.PostJSON(requestURL, nil)
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		if response.StatusCode != http.StatusOK {
			log.Fatalf("%s: %s", http.StatusText(response.StatusCode), body)
		}

		fmt.Println("Declined")
	},
}

func init() {
	rootCmd.AddCommand(declineCmd)
	declineCmd.Flags().StringVar(&Token, "token", "", "Token")
	declineCmd.MarkFlagRequired("token")
}
