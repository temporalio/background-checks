package mocks

import "github.com/temporalio/background-checks/internal"

var MotorVehicleIncidentSearchResults = map[internal.MotorVehicleIncidentSearchInput]internal.MotorVehicleIncidentSearchResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
