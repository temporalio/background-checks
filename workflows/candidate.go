package workflows

import (
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func Candidate(ctx workflow.Context, input types.CandidateInput) error {
	logger := workflow.GetLogger(ctx)

	checks := map[string]types.CandidateBackgroundCheckStatus{}

	err := workflow.SetQueryHandler(ctx, queries.CandidateBackgroundCheckList, func() ([]types.CandidateBackgroundCheckStatus, error) {
		result := make([]types.CandidateBackgroundCheckStatus, 0, len(checks))

		for _, check := range checks {
			result = append(result, check)
		}

		return result, nil
	})
	if err != nil {
		return err
	}

	s := workflow.NewSelector(ctx)

	createCh := workflow.GetSignalChannel(ctx, signals.CandidateBackgroundCheckStatus)
	s.AddReceive(createCh, func(c workflow.ReceiveChannel, more bool) {
		var bc types.CandidateBackgroundCheckStatus
		c.Receive(ctx, &bc)
		checks[bc.ID] = bc
	})

	consentCh := workflow.GetSignalChannel(ctx, signals.CandidateConsentFromUser)
	s.AddReceive(consentCh, func(c workflow.ReceiveChannel, more bool) {
		var consent types.CandidateConsentResponseFromUser
		c.Receive(ctx, &consent)

		f := workflow.SignalExternalWorkflow(
			ctx,
			consent.WorkflowID,
			consent.RunID,
			signals.CandidateConsentResponse,
			types.CandidateConsentResponse{
				Consent: consent.Consent,
			},
		)
		err := f.Get(ctx, nil)
		if err != nil {
			logger.Error("failed to send consent response from user: %v", err)
		}
	})

	for {
		s.Select(ctx)
	}
}
