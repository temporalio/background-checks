package workflows

import (
	"fmt"

	"github.com/temporalio/background-checks/internal"
	"go.temporal.io/sdk/workflow"
)

type BackgroundCheckStatus struct {
	AcceptResult                     internal.AcceptCheckResult
	ValidateSSNResult                internal.ValidateSSNResult
	FederalCriminalSearchResult      internal.FederalCriminalSearchResult
	StateCriminalSearchResult        internal.StateCriminalSearchResult
	MotorVehicleIncidentSearchResult internal.MotorVehicleIncidentSearchResult
}

func BackgroundCheck(ctx workflow.Context, email string, tier string) error {
	status := BackgroundCheckStatus{}

	logger := workflow.GetLogger(ctx)

	acceptWF := workflow.ExecuteChildWorkflow(ctx, AcceptCheck, internal.AcceptCheckInput{Email: email})
	err := acceptWF.Get(ctx, &status.AcceptResult)
	if err != nil {
		return err
	}

	if !status.AcceptResult.Accept {
		return nil
	}

	ssnInput := internal.ValidateSSNInput{
		FullName: status.AcceptResult.FullName,
		Address:  status.AcceptResult.Address,
		SSN:      status.AcceptResult.SSN,
	}
	ssnWF := workflow.ExecuteChildWorkflow(ctx, ValidateSSN, ssnInput)
	err = ssnWF.Get(ctx, &status.ValidateSSNResult)
	if err != nil {
		return err
	}

	if !status.ValidateSSNResult.Valid {
		return nil
	}

	s := workflow.NewSelector(ctx)

	federalCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		FederalCriminalSearch,
		internal.FederalCriminalSearchInput{FullName: status.AcceptResult.FullName, Address: status.AcceptResult.Address},
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
		internal.StateCriminalSearchInput{FullName: status.AcceptResult.FullName, Address: status.AcceptResult.Address},
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
		internal.MotorVehicleIncidentSearchInput{FullName: status.AcceptResult.FullName, Address: status.AcceptResult.Address},
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
