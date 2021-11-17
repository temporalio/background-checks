package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func waitForSubmission(ctx workflow.Context) types.Accept {
	var response types.Accept

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, signals.AcceptSubmission)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.AcceptSubmissionSignal
		c.Receive(ctx, &submission)

		response = submission.Accept
	})

	s.Select(ctx)

	return response
}

func Accept(ctx workflow.Context, input types.AcceptWorkflowInput) (types.AcceptWorkflowResult, error) {
	consent := waitForSubmission(ctx)

	return types.AcceptWorkflowResult{Accept: consent}, nil
}
