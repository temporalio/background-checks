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
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/background-checks/ui"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Run UI Server",
	Run: func(cmd *cobra.Command, args []string) {
		srv := &http.Server{
			Handler: ui.Router(),
			Addr:    "0.0.0.0:8083",
		}

		errCh := make(chan error, 1)
		go func() { errCh <- srv.ListenAndServe() }()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case <-sigCh:
			srv.Close()
		case err := <-errCh:
			log.Fatalf("error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
