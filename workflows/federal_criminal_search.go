package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func FederalCriminalSearch(ctx workflow.Context, input types.FederalCriminalSearchWorkflowInput) (types.FederalCriminalSearchWorkflowResult, error) {
	var result types.FederalCriminalSearchResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.FederalCriminalSearch, types.FederalCriminalSearchInput(input))

	err := f.Get(ctx, &result)
	return types.FederalCriminalSearchWorkflowResult(result), err
}
