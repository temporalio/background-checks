package workflows

import (
	"fmt"

	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

type BackgroundCheckStatus struct {
	Email                            string
	Tier                             string
	ConsentResult                    types.ConsentResult
	ValidateSSNResult                types.ValidateSSNResult
	FederalCriminalSearchResult      types.FederalCriminalSearchResult
	StateCriminalSearchResult        types.StateCriminalSearchResult
	MotorVehicleIncidentSearchResult types.MotorVehicleIncidentSearchResult
}

func BackgroundCheck(ctx workflow.Context, input types.BackgroundCheckInput) error {
	status := BackgroundCheckStatus{
		Email: input.Email,
		Tier:  input.Tier,
	}

	logger := workflow.GetLogger(ctx)

	acceptWF := workflow.ExecuteChildWorkflow(ctx, Consent, types.ConsentInput{Email: input.Email})
	err := acceptWF.Get(ctx, &status.ConsentResult)
	if err != nil {
		return err
	}

	if !status.ConsentResult.Accept {
		return nil
	}

	ssnInput := types.ValidateSSNInput{
		FullName: status.ConsentResult.FullName,
		Address:  status.ConsentResult.Address,
		SSN:      status.ConsentResult.SSN,
	}
	ssnWF := workflow.ExecuteChildWorkflow(ctx, ValidateSSN, ssnInput)
	err = ssnWF.Get(ctx, &status.ValidateSSNResult)
	if err != nil {
		return err
	}

	if !status.ValidateSSNResult.Valid {
		return nil
	}

	if input.Tier != "full" {
		return nil
	}

	s := workflow.NewSelector(ctx)

	federalCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		FederalCriminalSearch,
		types.FederalCriminalSearchInput{FullName: status.ConsentResult.FullName, Address: status.ConsentResult.Address},
	)
	s.AddFuture(federalCriminalSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.FederalCriminalSearchResult)
		if err != nil {
			logger.Error(fmt.Sprintf("federal criminal search: %v", err))
		}
	})

	stateCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		StateCriminalSearch,
		types.StateCriminalSearchInput{FullName: status.ConsentResult.FullName, Address: status.ConsentResult.Address},
	)
	s.AddFuture(stateCriminalSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.StateCriminalSearchResult)
		if err != nil {
			logger.Error(fmt.Sprintf("state criminal search: %v", err))
		}
	})

	motorVehicleIncidentSearch := workflow.ExecuteChildWorkflow(
		ctx,
		MotorVehicleIncidentSearch,
		types.MotorVehicleIncidentSearchInput{FullName: status.ConsentResult.FullName, Address: status.ConsentResult.Address},
	)
	s.AddFuture(motorVehicleIncidentSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.MotorVehicleIncidentSearchResult)
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
