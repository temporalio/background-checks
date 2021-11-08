package mocks

import "github.com/temporalio/background-checks/internal"

var FederalCriminalSearchResults = map[internal.FederalCriminalSearchInput]internal.FederalCriminalSearchResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
