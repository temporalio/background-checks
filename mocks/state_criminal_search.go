package mocks

import "github.com/temporalio/background-checks/types"

var StateCriminalSearchResults = map[types.StateCriminalSearchWorkflowInput]types.StateCriminalSearchWorkflowResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
