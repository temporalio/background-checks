package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func FederalCriminalSearch(ctx workflow.Context, input types.FederalCriminalSearchWorkflowInput) (types.FederalCriminalSearchWorkflowResult, error) {
	return mocks.FederalCriminalSearchResults[input], nil
}
