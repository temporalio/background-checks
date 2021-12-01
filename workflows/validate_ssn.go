package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func ValidateSSN(ctx workflow.Context, input types.ValidateSSNWorkflowInput) (types.ValidateSSNWorkflowResult, error) {
	var result types.SSNTraceResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SSNTrace, types.ValidateSSNWorkflowInput(input))

	err := f.Get(ctx, &result)
	return types.ValidateSSNWorkflowResult(result), err
}
