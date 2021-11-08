package workflows

import (
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/mocks"
	"go.temporal.io/sdk/workflow"
)

func MotorVehicleIncidentSearch(ctx workflow.Context, input internal.MotorVehicleIncidentSearchInput) (internal.MotorVehicleIncidentSearchResult, error) {
	return mocks.MotorVehicleIncidentSearchResults[input], nil
}
