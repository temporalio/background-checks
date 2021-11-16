package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func requestConsent(ctx workflow.Context, email string) error {
	f := workflow.SignalExternalWorkflow(
		ctx,
		mappings.CandidateWorkflowID(email),
		"",
		signals.ConsentRequest,
		types.ConsentRequestSignal{},
	)
	return f.Get(ctx, nil)
}

func waitForResponse(ctx workflow.Context) (types.ConsentResponseSignal, error) {
	var response types.ConsentResponseSignal

	s := workflow.NewSelector(ctx)

	consentCh := workflow.GetSignalChannel(ctx, signals.ConsentResponse)
	s.AddReceive(consentCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &response)
	})

	s.Select(ctx)

	return response, nil
}

func Consent(ctx workflow.Context, input types.ConsentWorkflowInput) (types.Consent, error) {
	result := types.Consent{}

	err := requestConsent(ctx, input.Email)
	if err != nil {
		return result, err
	}

	response, err := waitForResponse(ctx)
	if err != nil {
		return result, err
	}

	return response.Consent, nil
}
