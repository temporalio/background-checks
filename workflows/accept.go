package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func waitForSubmission(ctx workflow.Context) types.AcceptSubmission {
	var response types.AcceptSubmission

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, signals.AcceptSubmission)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.AcceptSubmissionSignal
		c.Receive(ctx, &submission)

		response = types.AcceptSubmission(submission)
	})

	s.Select(ctx)

	return response
}

func Accept(ctx workflow.Context, input types.AcceptWorkflowInput) (types.AcceptWorkflowResult, error) {
	submission := waitForSubmission(ctx)

	return types.AcceptWorkflowResult(submission), nil
}
