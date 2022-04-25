package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"go.temporal.io/sdk/workflow"
)

type SSNTraceWorkflowInput struct {
	FullName string
	SSN      string
}

type SSNTraceWorkflowResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

// @@@SNIPSTART background-checks-ssn-trace-workflow-definition

// SSNTrace is a Workflow Definition that calls for the execution of a single Activity.
// This is executed as a Child Workflow by the main Background Check.
func SSNTrace(ctx workflow.Context, input *SSNTraceWorkflowInput) (*SSNTraceWorkflowResult, error) {
	var result activities.SSNTraceResult

	ctx = workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteLocalActivity(ctx, a.SSNTrace, &activities.SSNTraceInput{
		FullName: input.FullName,
		SSN:      input.SSN,
	})

	err := f.Get(ctx, &result)
	r := SSNTraceWorkflowResult(result)
	return &r, err
}

// @@@SNIPEND
