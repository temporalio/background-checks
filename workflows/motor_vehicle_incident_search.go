package workflows

import (
	"time"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART background-checks-motor-vehicle-workflow-definition
func MotorVehicleIncidentSearch(ctx workflow.Context, input *types.MotorVehicleIncidentSearchWorkflowInput) (*types.MotorVehicleIncidentSearchWorkflowResult, error) {
	var result types.MotorVehicleIncidentSearchWorkflowResult

	name := input.FullName
	address := input.Address
	var motorvehicleIncidents []string

	activityInput := types.MotorVehicleIncidentSearchInput{
		FullName: name,
		Address:  address,
	}
	var activityResult types.MotorVehicleIncidentSearchResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	motorvehicleIncidentSearch := workflow.ExecuteActivity(ctx, a.MotorVehicleIncidentSearch, activityInput)

	err := motorvehicleIncidentSearch.Get(ctx, &activityResult)
	if err == nil {
		motorvehicleIncidents = append(motorvehicleIncidents, activityResult.MotorVehicleIncidents...)
	}
	result.MotorVehicleIncidents = motorvehicleIncidents

	r := types.MotorVehicleIncidentSearchWorkflowResult(result)
	return &r, nil
}

// @@@SNIPEND
