package workflows

import (
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func Researcher(ctx workflow.Context, input types.ResearcherInput) error {
	todos := map[string]types.ResearcherTodo{}

	err := workflow.SetQueryHandler(ctx, queries.ResearcherTodosList, func() ([]types.ResearcherTodo, error) {
		result := make([]types.ResearcherTodo, 0, len(todos))

		for _, todo := range todos {
			result = append(result, todo)
		}

		return result, nil
	})
	if err != nil {
		return err
	}

	s := workflow.NewSelector(ctx)

	createCh := workflow.GetSignalChannel(ctx, signals.ResearcherTodoCreate)
	s.AddReceive(createCh, func(c workflow.ReceiveChannel, more bool) {
		var todo types.ResearcherTodo
		c.Receive(ctx, &todo)
		todos[todo.Token] = todo
	})

	completeCh := workflow.GetSignalChannel(ctx, signals.ResearcherTodoComplete)
	s.AddReceive(completeCh, func(c workflow.ReceiveChannel, more bool) {
		var todo types.ResearcherTodo
		c.Receive(ctx, &todo)
		delete(todos, todo.Token)
	})

	for {
		s.Select(ctx)
	}
}
