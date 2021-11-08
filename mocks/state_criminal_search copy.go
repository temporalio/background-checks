package mocks

import "github.com/temporalio/background-checks/internal"

var StateCriminalSearchResults = map[internal.StateCriminalSearchInput]internal.StateCriminalSearchResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
