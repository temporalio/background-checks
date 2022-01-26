package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART background-checks-ssn-trace-workflow-definition
func SSNTrace(ctx workflow.Context, input *types.SSNTraceWorkflowInput) (*types.SSNTraceWorkflowResult, error) {
	var result types.SSNTraceResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SSNTrace, types.SSNTraceWorkflowInput(*input))

	err := f.Get(ctx, &result)
	r := types.SSNTraceWorkflowResult(result)
	return &r, err
}

// @@@SNIPEND
