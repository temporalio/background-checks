package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func SSNTrace(ctx workflow.Context, input types.SSNTraceWorkflowInput) (types.SSNTraceWorkflowResult, error) {
	var result types.SSNTraceResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SSNTrace, types.SSNTraceWorkflowInput(input))

	err := f.Get(ctx, &result)
	return types.SSNTraceWorkflowResult(result), err
}
