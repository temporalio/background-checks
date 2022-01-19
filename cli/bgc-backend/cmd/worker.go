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

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := temporal.NewClient(client.Options{
			MetricsHandler: tallyhandler.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
				ListenAddress: "0.0.0.0:8001",
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
		w.RegisterActivity(&activities.Activities{SMTPHost: "mailhog", SMTPPort: 1025})
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
		Prefix:         "",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)

	return scope
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
