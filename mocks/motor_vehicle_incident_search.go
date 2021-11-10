package mocks

import "github.com/temporalio/background-checks/types"

var MotorVehicleIncidentSearchResults = map[types.MotorVehicleIncidentSearchInput]types.MotorVehicleIncidentSearchResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
