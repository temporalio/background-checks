package workflows_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestEmploymentVerificationWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	details := workflows.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.OnActivity(a.SendEmploymentVerificationRequestEmail, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input *activities.SendEmploymentVerificationEmailInput) (*activities.SendEmploymentVerificationEmailResult, error) {
			return &activities.SendEmploymentVerificationEmailResult{}, nil
		},
	)

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				workflows.EmploymentVerificationSubmissionSignalName,
				workflows.EmploymentVerificationSubmissionSignal{EmploymentVerificationComplete: true, EmployerVerified: true},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.EmploymentVerification, &workflows.EmploymentVerificationWorkflowInput{CandidateDetails: details})

	var result workflows.EmploymentVerificationWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, workflows.EmploymentVerificationWorkflowResult{EmploymentVerificationComplete: true, EmployerVerified: true}, result)
}

func TestEmploymentVerificationWorkflowTimeout(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	details := workflows.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.OnActivity(a.SendEmploymentVerificationRequestEmail, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input *activities.SendEmploymentVerificationEmailInput) (*activities.SendEmploymentVerificationEmailResult, error) {
			return &activities.SendEmploymentVerificationEmailResult{}, nil
		},
	)

	env.ExecuteWorkflow(workflows.EmploymentVerification, &workflows.EmploymentVerificationWorkflowInput{CandidateDetails: details})

	var result workflows.EmploymentVerificationWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, workflows.EmploymentVerificationWorkflowResult{EmploymentVerificationComplete: false, EmployerVerified: false}, result)
}
