package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"go.temporal.io/sdk/workflow"
)

type MotorVehicleIncidentSearchWorkflowInput struct {
	FullName string
	Address  string
}

type MotorVehicleIncidentSearchWorkflowResult struct {
	LicenseValid          bool
	MotorVehicleIncidents []string
}

// @@@SNIPSTART background-checks-motor-vehicle-workflow-definition
func MotorVehicleIncidentSearch(ctx workflow.Context, input *MotorVehicleIncidentSearchWorkflowInput) (*MotorVehicleIncidentSearchWorkflowResult, error) {
	var result MotorVehicleIncidentSearchWorkflowResult

	name := input.FullName
	address := input.Address
	var motorvehicleIncidents []string

	activityInput := activities.MotorVehicleIncidentSearchInput{
		FullName: name,
		Address:  address,
	}
	var activityResult activities.MotorVehicleIncidentSearchResult

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	motorvehicleIncidentSearch := workflow.ExecuteActivity(ctx, a.MotorVehicleIncidentSearch, activityInput)

	err := motorvehicleIncidentSearch.Get(ctx, &activityResult)
	if err == nil {
		motorvehicleIncidents = append(motorvehicleIncidents, activityResult.MotorVehicleIncidents...)
	}
	result.MotorVehicleIncidents = motorvehicleIncidents

	r := MotorVehicleIncidentSearchWorkflowResult(result)
	return &r, nil
}

// @@@SNIPEND
