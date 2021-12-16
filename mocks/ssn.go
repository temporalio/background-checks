package mocks

import "github.com/temporalio/background-checks/types"

var SSNTraceWorkflowResult = map[types.SSNTraceWorkflowInput]types.SSNTraceWorkflowResult{
	{FullName: "John Smith", SSN: "111-11-1111"}:    {},
	{FullName: "Sally Jones", SSN: "123-45-6789"}:   {},
	{FullName: "Javier Bardem", SSN: "987-65-4321"}: {},
}

/*
{
	{Address: "123 Broadway", City: "New York", State: "NY", ZipCode: "10011"},
	{Address: "500 Market Street", City: "San Francisco", State: "CA", ZipCode: "94110"},
	{Address: "111 Dearborn Ave", City: "Detroit", State: "MI", ZipCode: "44014"},
} */
