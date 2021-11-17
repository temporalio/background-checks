package mocks

import "github.com/temporalio/background-checks/types"

var AcceptResultAccepted = types.Accept{
	Accept:   true,
	FullName: "John Smith",
	Address:  "1 Chestnut Avenue",
	SSN:      "111-11-1111",
}

var AcceptResultDeclined = map[string]types.Accept{
	"user1@example.com": AcceptResultAccepted,
}
