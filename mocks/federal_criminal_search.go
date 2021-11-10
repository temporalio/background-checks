package mocks

import "github.com/temporalio/background-checks/types"

var FederalCriminalSearchResults = map[types.FederalCriminalSearchInput]types.FederalCriminalSearchResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
