package workflows

import (
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/mocks"
	"go.temporal.io/sdk/workflow"
)

func AcceptCheck(ctx workflow.Context, input internal.AcceptCheckInput) (internal.AcceptCheckResult, error) {
	return mocks.AcceptCheckResults[input.Email], nil
}
