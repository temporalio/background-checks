package activities

import (
	"context"
	"encoding/base64"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

func Consent(ctx context.Context, input types.ConsentInput) error {
	info := activity.GetInfo(ctx)
	token := base64.StdEncoding.EncodeToString(info.TaskToken)

	c, err := client.NewClient(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.SignalWithStartWorkflow(
		ctx,
		mappings.CandidateWorkflowID(input.Email),
		signals.CandidateTodoCreate,
		types.CandidateTodo{
			Token: token,
		},
		client.StartWorkflowOptions{
			TaskQueue: config.TaskQueue,
		},
		"Candidate",
		types.CandidateInput{
			Email: input.Email,
		},
	)
	if err != nil {
		return err
	}

	return activity.ErrResultPending
}
