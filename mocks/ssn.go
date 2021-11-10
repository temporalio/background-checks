package mocks

import "github.com/temporalio/background-checks/types"

var ValidateSSNResults = map[types.ValidateSSNInput]types.ValidateSSNResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue", SSN: "111-11-1111"}: {Valid: true},
}
