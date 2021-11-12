package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func Consent(ctx workflow.Context, input types.ConsentInput) (types.ConsentResult, error) {
	info := workflow.GetInfo(ctx)

	f := workflow.SignalExternalWorkflow(
		ctx,
		mappings.CandidateWorkflowID(input.Email),
		"",
		signals.CandidateConsentRequest,
		types.CandidateConsentRequest{
			WorkflowID: info.WorkflowExecution.ID,
			RunID:      info.WorkflowExecution.RunID,
		},
	)
	err := f.Get(ctx, nil)
	if err != nil {
		return types.ConsentResult{}, err
	}

	var response types.CandidateConsentResponse

	s := workflow.NewSelector(ctx)

	consentCh := workflow.GetSignalChannel(ctx, signals.CandidateConsentFromUser)
	s.AddReceive(consentCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &response)
	})

	s.Select(ctx)

	return response.Consent, nil
}
