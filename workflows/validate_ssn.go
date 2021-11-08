package workflows

import (
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/mocks"
	"go.temporal.io/sdk/workflow"
)

func ValidateSSN(ctx workflow.Context, input internal.ValidateSSNInput) (internal.ValidateSSNResult, error) {
	return mocks.ValidateSSNResults[input], nil
}
