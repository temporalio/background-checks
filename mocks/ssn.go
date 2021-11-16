package mocks

import "github.com/temporalio/background-checks/types"

var ValidateSSNResults = map[types.ValidateSSNWorkflowInput]types.ValidateSSNWorkflowResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue", SSN: "111-11-1111"}: {Valid: true},
}
