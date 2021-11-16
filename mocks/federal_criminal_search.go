package mocks

import "github.com/temporalio/background-checks/types"

var FederalCriminalSearchResults = map[types.FederalCriminalSearchWorkflowInput]types.FederalCriminalSearchWorkflowResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
