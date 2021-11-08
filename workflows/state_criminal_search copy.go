package workflows

import (
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/mocks"
	"go.temporal.io/sdk/workflow"
)

func StateCriminalSearch(ctx workflow.Context, input internal.StateCriminalSearchInput) (internal.StateCriminalSearchResult, error) {
	return mocks.StateCriminalSearchResults[input], nil
}
