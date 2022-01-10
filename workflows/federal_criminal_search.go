package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART background-checks-federal-criminal-workflow-definition
func FederalCriminalSearch(ctx workflow.Context, input *types.FederalCriminalSearchWorkflowInput) (*types.FederalCriminalSearchWorkflowResult, error) {
	var result types.FederalCriminalSearchResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.FederalCriminalSearch, types.FederalCriminalSearchInput(*input))

	err := f.Get(ctx, &result)
	r := types.FederalCriminalSearchWorkflowResult(result)
	return &r, err
}

// @@@SNIPEND
