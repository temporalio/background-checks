package workflows

import (
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/types"
)

const (
	EmploymentVerificationSubmissionSignal = "employment-verification-submission"
	ResearchDeadline                       = time.Hour * 24 * 7
)

func chooseResearcher(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput) (string, error) {
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

func emailEmploymentVerificationRequest(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	evsend := workflow.ExecuteActivity(ctx, a.SendEmploymentVerificationRequestEmail, types.SendEmploymentVerificationEmailInput{
		Email:            email,
		CandidateDetails: input.CandidateDetails,
		Token:            workflows.TokenForWorkflow(ctx),
	})
	return evsend.Get(ctx, nil)
}

func waitForEmploymentVerificationSubmission(ctx workflow.Context) (types.EmploymentVerificationSubmission, error) {
	var response types.EmploymentVerificationSubmission
	var err error

	logger := workflow.GetLogger(ctx)

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, EmploymentVerificationSubmissionSignal)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.EmploymentVerificationSubmissionSignal
		c.Receive(ctx, &submission)

		logger.Info("Signal received: ", submission)

		response = types.EmploymentVerificationSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, ResearchDeadline), func(f workflow.Future) {
		err = f.Get(ctx, nil)

		// We should probably fail the (child) workflow here.
		response.EmploymentVerificationComplete = false
		response.EmployerVerified = false
	})

	s.Select(ctx)

	return response, err
}

// @@@SNIPSTART background-checks-employment-verification-workflow-definition
func EmploymentVerification(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput) (types.EmploymentVerificationWorkflowResult, error) {
	researcher, err := chooseResearcher(ctx, input)
	if err != nil {
		return types.EmploymentVerificationWorkflowResult{}, err
	}

	err = emailEmploymentVerificationRequest(ctx, input, researcher)
	if err != nil {
		return types.EmploymentVerificationWorkflowResult{}, err
	}
	submission, err := waitForEmploymentVerificationSubmission(ctx)

	return types.EmploymentVerificationWorkflowResult(submission), err
}

// @@@SNIPEND
