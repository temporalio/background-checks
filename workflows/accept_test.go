package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestReturnsAcceptWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := activities.Activities{SMTPStub: true}

	env.RegisterActivity(a.SendAcceptEmail)

	details := workflows.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				workflows.AcceptSubmissionSignalName,
				workflows.AcceptSubmissionSignal{Accepted: true, CandidateDetails: details},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.Accept, &workflows.AcceptWorkflowInput{})

	var result workflows.AcceptWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, workflows.AcceptWorkflowResult{Accepted: true, CandidateDetails: details}, result)
}

func TestReturnsAcceptWorkflowTimeout(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := activities.Activities{SMTPStub: true}

	env.RegisterActivity(a.SendAcceptEmail)

	env.ExecuteWorkflow(workflows.Accept, &workflows.AcceptWorkflowInput{})

	var result workflows.AcceptWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, workflows.AcceptWorkflowResult{Accepted: false, CandidateDetails: workflows.CandidateDetails{}}, result)
}
