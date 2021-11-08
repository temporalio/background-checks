package mocks

import "github.com/temporalio/background-checks/internal"

var AcceptCheckResultAccepted = internal.AcceptCheckResult{
	Accept:   true,
	FullName: "John Smith",
	Address:  "1 Chestnut Avenue",
	SSN:      "111-11-1111",
}

var AcceptCheckResults = map[string]internal.AcceptCheckResult{
	"user1@example.com": AcceptCheckResultAccepted,
}
