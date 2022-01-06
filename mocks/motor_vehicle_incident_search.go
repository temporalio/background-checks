package mocks

import "github.com/temporalio/background-checks/types"

var MotorVehicleIncidentSearchResults = map[types.MotorVehicleIncidentSearchWorkflowInput]types.MotorVehicleIncidentSearchWorkflowResult{
	{FullName: "John Smith", Address: "1 Chestnut Avenue"}: {},
}
