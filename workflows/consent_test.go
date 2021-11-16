package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestConsentWorkflowRequestsConsent(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	env.OnSignalExternalWorkflow(
		"default-test-namespace",
		mappings.CandidateWorkflowID("user@example.com"),
		"",
		signals.ConsentRequest,
		types.ConsentRequestSignal{},
	).Return(nil).Once()

	env.ExecuteWorkflow(workflows.Consent, types.ConsentWorkflowInput{Email: "user@example.com"})
}

func TestReturnsConsentResponse(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	consent := types.Consent{
		Consent:  true,
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.OnSignalExternalWorkflow(
		"default-test-namespace",
		mappings.CandidateWorkflowID("user@example.com"),
		"",
		signals.ConsentRequest,
		types.ConsentRequestSignal{},
	).Return(nil).Once()

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.ConsentResponse,
				types.ConsentResponseSignal{Consent: consent},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.Consent, types.ConsentWorkflowInput{Email: "user@example.com"})

	var result types.Consent
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, consent, result)
}
