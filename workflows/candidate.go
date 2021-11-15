package workflows

import (
	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func sendConsentResponse(ctx workflow.Context, email string, result types.ConsentResponse) error {
	f := workflow.SignalExternalWorkflow(
		ctx,
		mappings.ConsentWorkflowID(email),
		"",
		signals.ConsentResponse,
		types.ConsentResponse(result),
	)
	return f.Get(ctx, nil)
}

func Candidate(ctx workflow.Context, input types.CandidateInput) error {
	logger := workflow.GetLogger(ctx)

	email := input.Email

	var check types.CandidateBackgroundCheckStatus

	err := workflow.SetQueryHandler(ctx, queries.CandidateBackgroundCheckStatus, func() (types.CandidateBackgroundCheckStatus, error) {
		return check, nil
	})
	if err != nil {
		return err
	}

	s := workflow.NewSelector(ctx)

	createCh := workflow.GetSignalChannel(ctx, signals.BackgroundCheckStatus)
	s.AddReceive(createCh, func(c workflow.ReceiveChannel, more bool) {
		var bc types.CandidateBackgroundCheckStatus
		c.Receive(ctx, &bc)
		check = bc
	})

	consentRequestCh := workflow.GetSignalChannel(ctx, signals.ConsentRequest)
	s.AddReceive(consentRequestCh, func(c workflow.ReceiveChannel, more bool) {
		var r types.ConsentRequest
		c.Receive(ctx, &r)
		check.ConsentRequired = true
	})

	submissionCh := workflow.GetSignalChannel(ctx, signals.ConsentSubmission)
	s.AddReceive(submissionCh, func(c workflow.ReceiveChannel, more bool) {
		var submission types.ConsentSubmission
		c.Receive(ctx, &submission)

		err := sendConsentResponse(ctx, email, types.ConsentResponse(submission))
		if err != nil {
			logger.Error("failed to send consent response from user: %v", err)
		}
		check.ConsentRequired = false
	})

	for {
		s.Select(ctx)
	}
}
