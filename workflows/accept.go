package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

// @@@SNIPSTART background-checks-candidate-accept-email-candidate
func emailCandidate(ctx workflow.Context, input types.AcceptWorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendAcceptEmail, types.SendAcceptEmailInput(input))
	return f.Get(ctx, nil)
}
// @@@SNIPEND

// @@@SNIPSTART background-checks-candidate-accept-wait-for-submission
func waitForSubmission(ctx workflow.Context) types.AcceptSubmission {
	var response types.AcceptSubmission

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, signals.AcceptSubmission)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.AcceptSubmissionSignal
		c.Receive(ctx, &submission)

		response = types.AcceptSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, config.AcceptGracePeriod), func(f workflow.Future) {
		// Treat failure to accept in time as declining.
		response.Accepted = false
	})

	s.Select(ctx)

	return response
}
// @@@SNIPEND

// @@@SNIPSTART background-checks-candidate-accept-workflow-definition
func Accept(ctx workflow.Context, input types.AcceptWorkflowInput) (types.AcceptWorkflowResult, error) {
	err := emailCandidate(ctx, input)
	if err != nil {
		return types.AcceptWorkflowResult{}, err
	}

	submission := waitForSubmission(ctx)

	return types.AcceptWorkflowResult(submission), nil
}
// @@@SNIPEND
