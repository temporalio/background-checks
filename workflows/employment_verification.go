package workflows

import (
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/types"
)

const (
	EmploymentVerificationDetailsQuery     = "employment-verification-details"
	EmploymentVerificationSubmissionSignal = "employment-verification-submission"
	ResearchDeadline                       = time.Hour * 24 * 7
)

// chooseResearcher encapsulates the logic that randomly chooses a Researcher using a Side Effect.
func chooseResearcher(ctx workflow.Context, input *types.EmploymentVerificationWorkflowInput) (string, error) {
	researchers := []string{
		"researcher1@example.com",
		"researcher2@example.com",
		"researcher3@example.com",
	}

	// Here we just pick a random researcher.
	// In a real use case you may round-robin, decide based on price or current workload,
	// or fetch a researcher from a third party API.
	var researcher string
	r := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return researchers[rand.Intn(len(researchers))]
	})
	err := r.Get(&researcher)

	return researcher, err
}

// emailEmploymentVerificationRequest encapsulates the logic that calls for the execution an Activity.
func emailEmploymentVerificationRequest(ctx workflow.Context, input *types.EmploymentVerificationWorkflowInput, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	evsend := workflow.ExecuteActivity(ctx, a.SendEmploymentVerificationRequestEmail, types.SendEmploymentVerificationEmailInput{
		Email: email,
		Token: TokenForWorkflow(ctx),
	})
	return evsend.Get(ctx, nil)
}

// waitForEmploymentVerificationSubmission encapsulates the logic that waits on and handles a Signal.
func waitForEmploymentVerificationSubmission(ctx workflow.Context) (*types.EmploymentVerificationSubmission, error) {
	var response types.EmploymentVerificationSubmission
	var err error

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, EmploymentVerificationSubmissionSignal)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.EmploymentVerificationSubmissionSignal
		c.Receive(ctx, &submission)

		response = types.EmploymentVerificationSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, ResearchDeadline), func(f workflow.Future) {
		err = f.Get(ctx, nil)

		// We should probably fail the (child) workflow here.
		response.EmploymentVerificationComplete = false
		response.EmployerVerified = false
	})

	s.Select(ctx)

	return &response, err
}

// @@@SNIPSTART background-checks-employment-verification-workflow-definition

// EmploymentVerification is a Workflow Definition that calls for the execution of a Side Effect, and an Activity,
// but then waits on and handles a Signal. It is also capable of handling a Query to get Candidate Details.
// This is executed as a Child Workflow by the main Background Check.
func EmploymentVerification(ctx workflow.Context, input *types.EmploymentVerificationWorkflowInput) (*types.EmploymentVerificationWorkflowResult, error) {
	var result types.EmploymentVerificationWorkflowResult

	err := workflow.SetQueryHandler(ctx, EmploymentVerificationDetailsQuery, func() (types.CandidateDetails, error) {
		return input.CandidateDetails, nil
	})
	if err != nil {
		return &result, err
	}

	researcher, err := chooseResearcher(ctx, input)
	if err != nil {
		return &result, err
	}

	err = emailEmploymentVerificationRequest(ctx, input, researcher)
	if err != nil {
		return &result, err
	}
	submission, err := waitForEmploymentVerificationSubmission(ctx)

	result = types.EmploymentVerificationWorkflowResult(*submission)
	return &result, err
}

// @@@SNIPEND
