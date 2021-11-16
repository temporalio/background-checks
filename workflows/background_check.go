package workflows

import (
	"fmt"
	"time"

	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func createCandidateWorkflow(ctx workflow.Context, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, activities.CreateCandidateWorkflow, types.CandidateWorkflowInput{Email: email})
	return f.Get(ctx, nil)
}

func updateCandidateCheckStatus(ctx workflow.Context, email string, status string) error {
	f := workflow.SignalExternalWorkflow(
		ctx,
		mappings.CandidateWorkflowID(email),
		"",
		signals.BackgroundCheckStatus,
		types.BackgroundCheckStatusSignal{
			Status: status,
		},
	)
	return f.Get(ctx, nil)
}

func waitForConsent(ctx workflow.Context, email string) (types.Consent, error) {
	var r types.Consent

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: mappings.ConsentWorkflowID(email),
	})
	consentWF := workflow.ExecuteChildWorkflow(ctx, Consent, types.ConsentWorkflowInput{Email: email})
	err := consentWF.Get(ctx, &r)

	return r, err
}

func BackgroundCheck(ctx workflow.Context, input types.BackgroundCheckWorkflowInput) error {
	email := input.Email

	status := types.BackgroundCheckStatus{
		Email: email,
		Tier:  input.Package,
	}

	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, queries.BackgroundCheckStatus, func() (types.BackgroundCheckStatus, error) {
		return status, nil
	})
	if err != nil {
		return err
	}

	err = createCandidateWorkflow(ctx, email)
	if err != nil {
		return err
	}

	err = updateCandidateCheckStatus(ctx, email, "Consent Required")
	if err != nil {
		return err
	}

	c, err := waitForConsent(ctx, email)
	if err != nil {
		return err
	}

	status.Consent = c

	if !c.Consent {
		return updateCandidateCheckStatus(ctx, email, "Declined")
	}

	err = updateCandidateCheckStatus(ctx, email, "In Progress")
	if err != nil {
		return err
	}

	err = updateCandidateCheckStatus(ctx, email, "Running")
	if err != nil {
		return err
	}

	ssnInput := types.ValidateSSNWorkflowInput{
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

	if input.Package != "full" {
		return nil
	}

	s := workflow.NewSelector(ctx)

	federalCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		FederalCriminalSearch,
		types.FederalCriminalSearchWorkflowInput{FullName: status.Consent.FullName, Address: status.Consent.Address},
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
		types.StateCriminalSearchWorkflowInput{FullName: status.Consent.FullName, Address: status.Consent.Address},
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
		types.MotorVehicleIncidentSearchWorkflowInput{FullName: status.Consent.FullName, Address: status.Consent.Address},
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
