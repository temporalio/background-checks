package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func Consent(ctx workflow.Context, input types.ConsentInput) (types.ConsentResult, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
	})

	f := workflow.ExecuteActivity(ctx, activities.Consent, input)

	var result types.ConsentResult
	err := f.Get(ctx, &result)

	return result, err
}
