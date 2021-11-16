package activities

import (
	"context"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/types"
)

func CreateCandidateWorkflow(ctx context.Context, input types.CandidateWorkflowInput) error {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: config.TaskQueue,
			ID:        mappings.CandidateWorkflowID(input.Email),
		},
		"Candidate",
		types.CandidateWorkflowInput{
			Email: input.Email,
		},
	)

	return err
}
