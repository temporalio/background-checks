package mocks

import "github.com/temporalio/background-checks/types"

var ConsentResultConsented = types.ConsentResult{
	Accept:   true,
	FullName: "John Smith",
	Address:  "1 Chestnut Avenue",
	SSN:      "111-11-1111",
}

var ConsentResultDeclined = map[string]types.ConsentResult{
	"user1@example.com": ConsentResultConsented,
}
