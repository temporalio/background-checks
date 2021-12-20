package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART background-checks-state-criminal-workflow-definition
func StateCriminalSearch(ctx workflow.Context, input types.StateCriminalSearchWorkflowInput) (types.StateCriminalSearchWorkflowResult, error) {
	var result types.StateCriminalSearchWorkflowResult

	name := input.FullName
	knownaddresses := input.SSNTraceResult
	var crimes []string

	// s := workflow.NewSelector(ctx)

	for _, address := range knownaddresses {
		var activityResult types.StateCriminalSearchResult
		var activityInput types.StateCriminalSearchInput
		activityInput.FullName = name
		activityInput.Address = address

		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
		})

		statecheck := workflow.ExecuteActivity(ctx, a.StateCriminalSearch, types.StateCriminalSearchInput(activityInput))

		err := statecheck.Get(ctx, &activityResult)
		if err != nil {
			crimes = append(crimes, activityResult.Crimes...)
		}

	}
	result.Crimes = crimes
	var err error

	return types.StateCriminalSearchWorkflowResult(result), err

}
// @@@SNIPEND
