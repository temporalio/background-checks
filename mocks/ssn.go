package mocks

import "github.com/temporalio/background-checks/internal"

var ValidateSSNResults = map[internal.ValidateSSNInput]internal.ValidateSSNResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue", SSN: "111-11-1111"}: {Valid: true},
}
