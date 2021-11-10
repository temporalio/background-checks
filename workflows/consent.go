package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func Consent(ctx workflow.Context, input types.ConsentInput) (types.ConsentResult, error) {
	return mocks.ConsentResultConsented, nil
}
