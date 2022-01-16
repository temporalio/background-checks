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
	"time"

	"github.com/spf13/cobra"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	tallyhandler "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/temporal"
	"github.com/temporalio/background-checks/workflows"
)

const (
	SMTPServer = "lp-mailhog:1025"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := temporal.NewClient(client.Options{
			MetricsHandler: tallyhandler.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
				ListenAddress: "0.0.0.0:9090",
				TimerType:     "histogram",
			})),
		})
		if err != nil {
			log.Fatalf("client error: %v", err)
		}
		defer c.Close()

		w := worker.New(c, "background-checks-main", worker.Options{})

		w.RegisterWorkflow(workflows.BackgroundCheck)
		w.RegisterWorkflow(workflows.Accept)
		w.RegisterWorkflow(workflows.EmploymentVerification)
		w.RegisterActivity(&activities.Activities{SMTPHost: "lp-mailhog", SMTPPort: 1025})
		w.RegisterWorkflow(workflows.SSNTrace)
		w.RegisterWorkflow(workflows.FederalCriminalSearch)
		w.RegisterWorkflow(workflows.StateCriminalSearch)
		w.RegisterWorkflow(workflows.MotorVehicleIncidentSearch)

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("worker exited: %v", err)
		}
	},
}

func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter: reporter,
		Separator:      prometheus.DefaultSeparator,
		Prefix:         "lp_background_checks",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)

	return scope
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
