package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"go.temporal.io/sdk/workflow"
)

type FederalCriminalSearchWorkflowInput struct {
	FullName       string
	KnownAddresses []string
}

type FederalCriminalSearchWorkflowResult struct {
	Crimes []string
}

// @@@SNIPSTART background-checks-federal-criminal-workflow-definition
func FederalCriminalSearch(ctx workflow.Context, input *FederalCriminalSearchWorkflowInput) (*FederalCriminalSearchWorkflowResult, error) {
	var result activities.FederalCriminalSearchResult

	name := input.FullName
	var address string
	if len(input.KnownAddresses) > 0 {
		address = input.KnownAddresses[0]
	}
	var crimes []string

	activityInput := activities.FederalCriminalSearchInput{
		FullName: name,
		Address:  address,
	}
	var activityResult activities.FederalCriminalSearchResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	federalcheck := workflow.ExecuteActivity(ctx, a.FederalCriminalSearch, activityInput)

	err := federalcheck.Get(ctx, &activityResult)
	if err == nil {
		crimes = append(crimes, activityResult.Crimes...)
	}
	result.Crimes = crimes

	r := FederalCriminalSearchWorkflowResult(result)
	return &r, nil
}

// @@@SNIPEND
