package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func FederalCriminalSearch(ctx workflow.Context, input types.FederalCriminalSearchInput) (types.FederalCriminalSearchResult, error) {
	return mocks.FederalCriminalSearchResults[input], nil
}
