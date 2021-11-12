package workflows

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/testsuite"

	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_CandidateBackgroundCheckNeedsConsent() {
	env := s.env

	// Emulate SignalWithStart
	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.CandidateBackgroundCheckStatus,
				types.CandidateBackgroundCheckStatus{
					ID:     "test",
					Status: "Awaiting Consent",
				},
			)
		},
		0,
	)

	env.ExecuteWorkflow(Candidate, types.CandidateInput{Email: "user@example.com"})

	v, err := s.env.QueryWorkflow(queries.CandidateBackgroundCheckList, nil)
	s.NoError(err)

	var list []types.CandidateBackgroundCheckStatus
	err = v.Get(&list)
	s.NoError(err)

	s.Equal([]types.CandidateBackgroundCheckStatus{
		{ID: "test", Status: "Awaiting Consent"},
	}, list)
}

func (s *UnitTestSuite) Test_CandidateProvidesConsent() {
	env := s.env

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.CandidateConsentRequest,
				types.CandidateConsentRequest{
					WorkflowID: "requestor-workflow-id",
					RunID:      "requestor-run-id",
				},
			)
		},
		0,
	)

	// Candidate sees consent is required and provides consent via CLI
	consent := types.ConsentResult{
		Consent:  true,
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.CandidateConsentFromUser,
				types.CandidateConsentResponseFromUser{
					WorkflowID: "requestor-workflow-id",
					RunID:      "requestor-run-id",
					Consent:    consent,
				},
			)
		},
		1,
	)

	env.OnSignalExternalWorkflow(
		"default-test-namespace",
		"requestor-workflow-id",
		"requestor-run-id",
		signals.CandidateConsentResponse,
		types.CandidateConsentResponse{
			Consent: consent,
		},
	).Return(nil).Once()

	env.ExecuteWorkflow(Candidate, types.CandidateInput{Email: "user@example.com"})
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
