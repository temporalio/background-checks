package workflows

import (
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/mocks"
	"go.temporal.io/sdk/workflow"
)

func FederalCriminalSearch(ctx workflow.Context, input internal.FederalCriminalSearchInput) (internal.FederalCriminalSearchResult, error) {
	return mocks.FederalCriminalSearchResults[input], nil
}
