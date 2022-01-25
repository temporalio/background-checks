package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART background-checks-federal-criminal-workflow-definition
func FederalCriminalSearch(ctx workflow.Context, input *types.FederalCriminalSearchWorkflowInput) (*types.FederalCriminalSearchWorkflowResult, error) {
	var result types.FederalCriminalSearchResult

	name := input.FullName
	var address string
	if len(input.KnownAddresses) > 0 {
		address = input.KnownAddresses[0]
	}
	var crimes []string

	activityInput := types.FederalCriminalSearchInput{
		FullName: name,
		Address:  address,
	}
	var activityResult types.FederalCriminalSearchResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	federalcheck := workflow.ExecuteActivity(ctx, a.FederalCriminalSearch, activityInput)

	err := federalcheck.Get(ctx, &activityResult)
	if err == nil {
		crimes = append(crimes, activityResult.Crimes...)
	}
	result.Crimes = crimes

	r := types.FederalCriminalSearchWorkflowResult(result)
	return &r, nil
}

// @@@SNIPEND
