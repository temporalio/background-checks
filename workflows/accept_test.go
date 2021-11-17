package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestReturnsAcceptResponse(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	accept := types.Accept{
		Accept:   true,
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.AcceptSubmission,
				types.AcceptSubmissionSignal{Accept: accept},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.Accept, types.AcceptWorkflowInput{})

	var result types.AcceptWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, types.AcceptWorkflowResult{Accept: accept}, result)
}
