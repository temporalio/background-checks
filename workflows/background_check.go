package workflows

import (
	"fmt"

	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/workflow"
)

func updateStatus(ctx workflow.Context, status types.BackgroundCheckStatus) error {
	return workflow.UpsertSearchAttributes(
		ctx,
		map[string]interface{}{
			"BackgroundCheckStatus": status.String(),
		},
	)
}

func waitForAccept(ctx workflow.Context, email string) (types.AcceptSubmission, error) {
	var r types.AcceptSubmission

	checkID := workflow.GetInfo(ctx).WorkflowExecution.RunID

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: mappings.AcceptWorkflowID(checkID),
	})
	consentWF := workflow.ExecuteChildWorkflow(ctx, Accept, types.AcceptWorkflowInput{
		Email:   email,
		CheckID: checkID,
	})
	err := consentWF.Get(ctx, &r)

	return r, err
}

func waitForEmploymentVerification(ctx workflow.Context, candidate types.CandidateDetails) (types.EmploymentVerificationWorkflowResult, error) {
	var r types.EmploymentVerificationWorkflowResult

	checkID := workflow.GetInfo(ctx).WorkflowExecution.RunID

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: mappings.EmploymentVerificationWorkflowID(checkID),
	})
	employmentVerificationWF := workflow.ExecuteChildWorkflow(ctx, EmploymentVerification, types.EmploymentVerificationWorkflowInput{
		CandidateDetails: candidate,
		CheckID:          checkID,
	})
	err := employmentVerificationWF.Get(ctx, &r)

	return r, err
}

func BackgroundCheck(ctx workflow.Context, input types.BackgroundCheckWorkflowInput) error {
	email := input.Email

	status := types.BackgroundCheckState{
		Email: email,
		Tier:  input.Package,
	}

	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, queries.BackgroundCheckStatus, func() (types.BackgroundCheckState, error) {
		return status, nil
	})
	if err != nil {
		return err
	}

	err = updateStatus(ctx, types.BackgroundCheckStatusPendingAccept)
	if err != nil {
		return err
	}

	c, err := waitForAccept(ctx, input.Email)
	if err != nil {
		return err
	}

	candidate := c.CandidateDetails
	status.CandidateDetails = candidate

	if !c.Accepted {
		return updateStatus(ctx, types.BackgroundCheckStatusDeclined)
	}

	err = updateStatus(ctx, types.BackgroundCheckStatusRunning)
	if err != nil {
		return err
	}

	ssnTrace := workflow.ExecuteChildWorkflow(
		ctx, ValidateSSN,
		types.ValidateSSNWorkflowInput{FullName: candidate.FullName, SSN: candidate.SSN},
	)

	err = ssnTrace.Get(ctx, &status.ValidateSSN)
	if err != nil {
		return err
	}

	s := workflow.NewSelector(ctx)

	federalCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		FederalCriminalSearch,
		types.FederalCriminalSearchWorkflowInput{FullName: candidate.FullName, Address: candidate.Address},
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
		types.StateCriminalSearchWorkflowInput{FullName: candidate.FullName, Address: candidate.Address},
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
		types.MotorVehicleIncidentSearchWorkflowInput{FullName: candidate.FullName, Address: candidate.Address},
	)
	s.AddFuture(motorVehicleIncidentSearch, func(f workflow.Future) {
		err := f.Get(ctx, &status.MotorVehicleIncidentSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("motor vehicle incident search: %v", err))
		}
	})

	// Employment Verification

	ev, err := waitForEmploymentVerification(ctx, candidate)
	if err != nil {
		return err
	}

	if ev.EmployerVerificationComplete {
		candidate.EmployerVerified = ev.CandidateDetails.EmployerVerified
	}

	checks := []workflow.ChildWorkflowFuture{
		federalCriminalSearch,
		stateCriminalSearch,
		motorVehicleIncidentSearch,
	}

	for i := 0; i < len(checks); i++ {
		s.Select(ctx)
	}

	return updateStatus(ctx, types.BackgroundCheckStatusCompleted)
}
