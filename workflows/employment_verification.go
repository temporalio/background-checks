package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func emailEmploymentVerificationRequest(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	evsend := workflow.ExecuteActivity(ctx, a.SendEmploymentVerificationRequestEmail, types.SendEmploymentVerificationEmailInput(input))
	return evsend.Get(ctx, nil)
}

func waitForEmploymentVerificationSubmission(ctx workflow.Context) types.EmploymentVerificationSubmission {
	var response types.EmploymentVerificationSubmission

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, signals.EmploymentVerificationSubmission)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.EmploymentVerificationSubmissionSignal
		c.Receive(ctx, &submission)

		response = types.EmploymentVerificationSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, config.AcceptGracePeriod), func(f workflow.Future) {
		// Treat failure to accept in time as declining.
		response.EmployerVerificationComplete = false
	})

	s.Select(ctx)

	return response
}

func EmploymentVerification(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput) (types.EmploymentVerificationWorkflowResult, error) {
	err := emailEmploymentVerificationRequest(ctx, input)
	if err != nil {
		return types.EmploymentVerificationWorkflowResult{}, err
	}
	submission := waitForEmploymentVerificationSubmission(ctx)

	return types.EmploymentVerificationWorkflowResult(submission), nil
}
