package workflows

import (
	"github.com/temporalio/background-checks/mocks"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func MotorVehicleIncidentSearch(ctx workflow.Context, input types.MotorVehicleIncidentSearchInput) (types.MotorVehicleIncidentSearchResult, error) {
	return mocks.MotorVehicleIncidentSearchResults[input], nil
}
