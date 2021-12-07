package workflows

import (
	"fmt"
	"time"

	"github.com/temporalio/background-checks/config"
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

	checkID := workflow.GetInfo(ctx).WorkflowExecution.RunID

	state := types.BackgroundCheckState{
		Email: email,
		Tier:  input.Package,
	}

	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, queries.BackgroundCheckStatus, func() (types.BackgroundCheckState, error) {
		return state, nil
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

	state.CandidateDetails = c.CandidateDetails
	state.Accepted = c.Accepted

	if !c.Accepted {
		return updateStatus(ctx, types.BackgroundCheckStatusDeclined)
	}

	err = updateStatus(ctx, types.BackgroundCheckStatusRunning)
	if err != nil {
		return err
	}

	ssnTrace := workflow.ExecuteChildWorkflow(
		ctx,
		ValidateSSN,
		types.ValidateSSNWorkflowInput{FullName: state.CandidateDetails.FullName, SSN: state.CandidateDetails.SSN},
	)

	err = ssnTrace.Get(ctx, &state.ValidateSSN)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("ssn trace: %v", state.ValidateSSN))

	s := workflow.NewSelector(ctx)

	federalCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		FederalCriminalSearch,
		types.FederalCriminalSearchWorkflowInput{FullName: state.CandidateDetails.FullName, Address: state.CandidateDetails.Address},
	)
	s.AddFuture(federalCriminalSearch, func(f workflow.Future) {
		err := f.Get(ctx, &state.FederalCriminalSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("federal criminal search: %v", err))
		}
		logger.Info(fmt.Sprintf("Federal Search: %v", state.FederalCriminalSearch))
	})

	/* State check will iterate over array of Known Addresses
	 */

	stateCriminalSearch := workflow.ExecuteChildWorkflow(
		ctx,
		StateCriminalSearch,
		types.StateCriminalSearchWorkflowInput{FullName: state.CandidateDetails.FullName, SSNTraceResult: state.ValidateSSN.KnownAddresses},
	)
	s.AddFuture(stateCriminalSearch, func(f workflow.Future) {
		err := f.Get(ctx, &state.StateCriminalSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("state criminal search: %v", err))
		}
		logger.Info(fmt.Sprintf("State Search: %v", state.StateCriminalSearch))
	})

	motorVehicleIncidentSearch := workflow.ExecuteChildWorkflow(
		ctx,
		MotorVehicleIncidentSearch,
		types.MotorVehicleIncidentSearchWorkflowInput{FullName: state.CandidateDetails.FullName, Address: state.CandidateDetails.Address},
	)
	s.AddFuture(motorVehicleIncidentSearch, func(f workflow.Future) {
		err := f.Get(ctx, &state.MotorVehicleIncidentSearch)
		if err != nil {
			logger.Error(fmt.Sprintf("motor vehicle incident search: %v", err))
		}
		logger.Info(fmt.Sprintf("Motor Vehicle Search: %v", state.MotorVehicleIncidentSearch))
	})

	// Employment Verification

	employmentVerification := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: mappings.EmploymentVerificationWorkflowID(checkID),
		}),
		EmploymentVerification,
		types.EmploymentVerificationWorkflowInput{CandidateDetails: state.CandidateDetails, CheckID: checkID},
	)
	s.AddFuture(employmentVerification, func(f workflow.Future) {
		err := f.Get(ctx, &state.EmploymentVerification)
		if err != nil {
			logger.Error(fmt.Sprintf("employment verification: %v", err))
		}
		logger.Info(fmt.Sprintf("Employment Verification: %v", state.EmploymentVerification))
	})

	checks := []workflow.ChildWorkflowFuture{
		federalCriminalSearch,
		stateCriminalSearch,
		motorVehicleIncidentSearch,
		employmentVerification,
	}

	for i := 0; i < len(checks); i++ {
		s.Select(ctx)
	}

	updateStatus(ctx, types.BackgroundCheckStatusCompleted)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendReportEmail, types.SendReportEmailInput{Email: config.HiringManagerEmail, State: state})
	return f.Get(ctx, nil)
}
