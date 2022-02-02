package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"go.temporal.io/sdk/workflow"
)

type StateCriminalSearchWorkflowInput struct {
	FullName       string
	KnownAddresses []string
}

type StateCriminalSearchWorkflowResult struct {
	Crimes []string
}

// @@@SNIPSTART background-checks-state-criminal-workflow-definition

// StateCriminalSearch is a Workflow Definition that calls for the execution an Activity for
// each address associated with the Candidate.
// This is executed as a Child Workflow by the main Background Check.
func StateCriminalSearch(ctx workflow.Context, input *StateCriminalSearchWorkflowInput) (*StateCriminalSearchWorkflowResult, error) {
	var result StateCriminalSearchWorkflowResult

	name := input.FullName
	knownaddresses := input.KnownAddresses
	var crimes []string

	for _, address := range knownaddresses {
		activityInput := activities.StateCriminalSearchInput{
			FullName: name,
			Address:  address,
		}
		var activityResult activities.StateCriminalSearchResult

		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
		})

		statecheck := workflow.ExecuteActivity(ctx, a.StateCriminalSearch, activityInput)

		err := statecheck.Get(ctx, &activityResult)
		if err == nil {
			crimes = append(crimes, activityResult.Crimes...)
		}
	}
	result.Crimes = crimes

	r := StateCriminalSearchWorkflowResult(result)
	return &r, nil
}

// @@@SNIPEND
