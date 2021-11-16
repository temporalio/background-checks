package workflows_test

import (
	"testing"

	"go.temporal.io/sdk/testsuite"

	"github.com/stretchr/testify/assert"

	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
)

func Test_CandidateBackgroundCheckNeedsConsent(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.BackgroundCheckStatus,
				types.BackgroundCheckStatusSignal{
					Status:          "Consent Required",
					ConsentRequired: true,
				},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.Candidate, types.CandidateWorkflowInput{Email: "user@example.com"})

	v, err := env.QueryWorkflow(queries.CandidateBackgroundCheckStatus, nil)
	assert.NoError(t, err)

	var check types.BackgroundCheckStatusSignal
	err = v.Get(&check)
	assert.NoError(t, err)

	assert.Equal(t,
		types.BackgroundCheckStatusSignal{Status: "Consent Required", ConsentRequired: true},
		check,
	)
}

func Test_CandidateProvidesConsent(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.BackgroundCheckStatus,
				types.BackgroundCheckStatusSignal{
					Status:          "Consent Required",
					ConsentRequired: true,
				},
			)
		},
		0,
	)

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.ConsentRequest,
				types.ConsentRequestSignal{},
			)
		},
		1,
	)

	// Candidate sees consent is required and provides consent via CLI
	consent := types.Consent{
		Consent:  true,
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.ConsentSubmission,
				types.ConsentSubmissionSignal{
					Consent: consent,
				},
			)
		},
		2,
	)

	env.OnSignalExternalWorkflow(
		"default-test-namespace",
		mappings.ConsentWorkflowID("user@example.com"),
		"",
		signals.ConsentResponse,
		types.ConsentResponseSignal{
			Consent: consent,
		},
	).Return(nil).Once()

	env.ExecuteWorkflow(workflows.Candidate, types.CandidateWorkflowInput{Email: "user@example.com"})
}
