package workflows

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/testsuite"

	"github.com/temporalio/background-checks/mappings"
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

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.BackgroundCheckStatus,
				types.CandidateBackgroundCheckStatus{
					Status:          "Consent Required",
					ConsentRequired: true,
				},
			)
		},
		0,
	)

	env.ExecuteWorkflow(Candidate, types.CandidateInput{Email: "user@example.com"})

	v, err := s.env.QueryWorkflow(queries.CandidateBackgroundCheckStatus, nil)
	s.NoError(err)

	var check types.CandidateBackgroundCheckStatus
	err = v.Get(&check)
	s.NoError(err)

	s.Equal(
		types.CandidateBackgroundCheckStatus{Status: "Consent Required", ConsentRequired: true},
		check,
	)
}

func (s *UnitTestSuite) Test_CandidateProvidesConsent() {
	env := s.env

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				signals.BackgroundCheckStatus,
				types.CandidateBackgroundCheckStatus{
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
				types.ConsentRequest{},
			)
		},
		1,
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
				signals.ConsentSubmission,
				types.ConsentSubmission{
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
		types.ConsentResponse{
			Consent: consent,
		},
	).Return(nil).Once()

	env.ExecuteWorkflow(Candidate, types.CandidateInput{Email: "user@example.com"})
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
