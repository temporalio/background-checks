package workflows

import (
	"fmt"

	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func BackgroundCheck(ctx workflow.Context, input types.BackgroundCheckInput) error {
	status := types.BackgroundCheckStatus{
		Email: input.Email,
		Tier:  input.Tier,
	}

	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, queries.BackgroundCheckStatus, func() (types.BackgroundCheckStatus, error) {
		return status, nil
	})
	if err != nil {
		return err
	}

	acceptWF := workflow.ExecuteChildWorkflow(ctx, Consent, types.ConsentInput{Email: input.Email})
	err = acceptWF.Get(ctx, &status.Consent)
	if err != nil {
		return err
	}

	if !status.Consent.Consent {
		return nil
	}

	ssnInput := types.ValidateSSNInput{
		FullName: status.Consent.FullName,
		Address:  status.Consent.Address,
		SSN:      status.Consent.SSN,
	}
	ssnWF := workflow.ExecuteChildWorkflow(ctx, ValidateSSN, ssnInput)
	err = ssnWF.Get(ctx, &status.Validate)
	if err != nil {
		return err
	}

	if !status.Validate.Valid {
		return nil
	}

	if input.Tier != "full" {
		return nil
	}

	s := workflow.NewSelector(ctx)

	federalCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		FederalCriminalSearch,
		types.FederalCriminalSearchInput{FullName: status.Consent.FullName, Address: status.Consent.Address},
	)
	s.AddFuture(federalCriminalSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.FederalCriminalSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("federal criminal search: %v", err))
		}
	})

	stateCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		StateCriminalSearch,
		types.StateCriminalSearchInput{FullName: status.Consent.FullName, Address: status.Consent.Address},
	)
	s.AddFuture(stateCriminalSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.StateCriminalSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("state criminal search: %v", err))
		}
	})

	motorVehicleIncidentSearch := workflow.ExecuteChildWorkflow(
		ctx,
		MotorVehicleIncidentSearch,
		types.MotorVehicleIncidentSearchInput{FullName: status.Consent.FullName, Address: status.Consent.Address},
	)
	s.AddFuture(motorVehicleIncidentSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.MotorVehicleIncidentSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("motor vehicle incident search: %v", err))
		}
	})

	checks := []workflow.ChildWorkflowFuture{
		federalCriminalSearch,
		stateCriminalSearch,
		motorVehicleIncidentSearch,
	}

	for i := 0; i < len(checks); i++ {
		s.Select(ctx)
	}

	return nil
}
