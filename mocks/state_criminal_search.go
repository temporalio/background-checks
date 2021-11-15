package mocks

import "github.com/temporalio/background-checks/types"

var StateCriminalSearchResults = map[types.StateCriminalSearchInput]types.StateCriminalSearchResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
