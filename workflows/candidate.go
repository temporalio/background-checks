package workflows

import (
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func Candidate(ctx workflow.Context, input types.CandidateInput) error {
	todos := map[string]types.CandidateTodo{}

	err := workflow.SetQueryHandler(ctx, queries.CandidateTodosList, func() ([]types.CandidateTodo, error) {
		result := make([]types.CandidateTodo, 0, len(todos))

		for _, todo := range todos {
			result = append(result, todo)
		}

		return result, nil
	})
	if err != nil {
		return err
	}

	s := workflow.NewSelector(ctx)

	createCh := workflow.GetSignalChannel(ctx, signals.CandidateTodoCreate)
	s.AddReceive(createCh, func(c workflow.ReceiveChannel, more bool) {
		var todo types.CandidateTodo
		c.Receive(ctx, &todo)
		todos[todo.Token] = todo
	})

	completeCh := workflow.GetSignalChannel(ctx, signals.CandidateTodoComplete)
	s.AddReceive(completeCh, func(c workflow.ReceiveChannel, more bool) {
		var todo types.CandidateTodo
		c.Receive(ctx, &todo)
		delete(todos, todo.Token)
	})

	for {
		s.Select(ctx)
	}
}
