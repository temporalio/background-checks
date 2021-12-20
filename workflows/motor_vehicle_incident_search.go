package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART background-checks-motor-vehicle-workflow-definition
func MotorVehicleIncidentSearch(ctx workflow.Context, input types.MotorVehicleIncidentSearchWorkflowInput) (types.MotorVehicleIncidentSearchWorkflowResult, error) {
	return mocks.MotorVehicleIncidentSearchResults[input], nil
}
// @@@SNIPEND
