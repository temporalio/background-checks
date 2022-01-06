package workflows_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestReturnsAcceptResponse(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	details := types.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.OnActivity(a.SendAcceptEmail, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input types.SendAcceptEmailInput) (types.SendAcceptEmailResult, error) {
			return types.SendAcceptEmailResult{}, nil
		},
	)

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.AcceptSubmission,
				types.AcceptSubmissionSignal{Accepted: true, CandidateDetails: details},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.Accept, types.AcceptWorkflowInput{})

	var result types.AcceptWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, types.AcceptWorkflowResult{Accepted: true, CandidateDetails: details}, result)
}
