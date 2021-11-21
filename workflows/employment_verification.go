package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func emailEmploymentVerificationRequest(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	evsend := workflow.ExecuteActivity(ctx, activities.SendEmploymentVerificationRequestEmail, types.SendEmploymentVerificationEmailInput(input))
	return evsend.Get(ctx, nil)
}

func waitForEmploymentVerificationSubmission(ctx workflow.Context) types.EmploymentVerificationWorkflowResult {
	var response types.EmploymentVerificationWorkflowResult

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, signals.EmploymentVerificationSubmission)

	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.EmploymentVerificationSubmission
		c.Receive(ctx, &submission)

		response = types.EmploymentVerificationWorkflowResult(submission)
	})

	s.Select(ctx)

	return response
}

func EmploymentVerification(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput) (types.EmploymentVerificationWorkflowResult, error) {
	submission := waitForEmploymentVerificationSubmission(ctx)

	return types.EmploymentVerificationWorkflowResult(submission), nil
}
