package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"go.temporal.io/sdk/workflow"
)

const (
	AcceptSubmissionSignalName = "accept-submission"
	AcceptGracePeriod          = time.Hour * 24 * 7
)

func emailCandidate(ctx workflow.Context, input *AcceptWorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	i := activities.SendAcceptEmailInput{
		Email: input.Email,
		Token: TokenForWorkflow(ctx),
	}
	f := workflow.ExecuteActivity(ctx, a.SendAcceptEmail, i)
	return f.Get(ctx, nil)
}

func waitForSubmission(ctx workflow.Context) (*AcceptSubmission, error) {
	var response AcceptSubmission
	var err error

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, AcceptSubmissionSignalName)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission AcceptSubmissionSignal
		c.Receive(ctx, &submission)

		response = AcceptSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, AcceptGracePeriod), func(f workflow.Future) {
		err = f.Get(ctx, nil)

		// Treat failure to accept in time as declining.
		response.Accepted = false
	})

	s.Select(ctx)

	return &response, err
}

type AcceptWorkflowInput struct {
	Email string
}

type AcceptWorkflowResult struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

// @@@SNIPSTART background-checks-accept-workflow-definition
func Accept(ctx workflow.Context, input *AcceptWorkflowInput) (*AcceptWorkflowResult, error) {
	err := emailCandidate(ctx, input)
	if err != nil {
		return &AcceptWorkflowResult{}, err
	}

	submission, err := waitForSubmission(ctx)

	result := AcceptWorkflowResult(*submission)
	return &result, err
}

// @@@SNIPEND
