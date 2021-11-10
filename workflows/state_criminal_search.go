package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func StateCriminalSearch(ctx workflow.Context, input types.StateCriminalSearchInput) (types.StateCriminalSearchResult, error) {
	return mocks.StateCriminalSearchResults[input], nil
}
