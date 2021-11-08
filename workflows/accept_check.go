package workflows

import (
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/mocks"
	"go.temporal.io/sdk/workflow"
)

func AcceptCheck(ctx workflow.Context, email string) (internal.AcceptCheckResult, error) {
	return mocks.AcceptCheckResults[email], nil
}
