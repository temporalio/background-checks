package workflows

import (
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
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

	researcher := researchers[rand.Intn(len(researchers))]

	return researcher, nil
}

func emailEmploymentVerificationRequest(ctx workflow.Context, input types.EmploymentVerificationWorkflowInput, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	evsend := workflow.ExecuteActivity(ctx, a.SendEmploymentVerificationRequestEmail, types.SendEmploymentVerificationEmailInput{
		Email:            email,
		CheckID:          input.CheckID,
		CandidateDetails: input.CandidateDetails,
	})
	return evsend.Get(ctx, nil)
}

func waitForEmploymentVerificationSubmission(ctx workflow.Context) types.EmploymentVerificationSubmission {
	var response types.EmploymentVerificationSubmission
	logger := workflow.GetLogger(ctx)

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, signals.EmploymentVerificationSubmission)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.EmploymentVerificationSubmissionSignal
		c.Receive(ctx, &submission)

		logger.Info("Signal received: ", submission)

		response = types.EmploymentVerificationSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, config.ResearchDeadline), func(f workflow.Future) {
		// We should probably fail the (child) workflow here.
		response.EmploymentVerificationComplete = false
		response.EmployerVerified = false
	})

	s.Select(ctx)

	return response
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
	submission := waitForEmploymentVerificationSubmission(ctx)

	return types.EmploymentVerificationWorkflowResult(submission), nil
}
// @@@SNIPEND
