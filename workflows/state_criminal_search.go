package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func StateCriminalSearch(ctx workflow.Context, input types.StateCriminalSearchWorkflowInput) (types.StateCriminalSearchWorkflowResult, error) {
	return mocks.StateCriminalSearchResults[input], nil
}
