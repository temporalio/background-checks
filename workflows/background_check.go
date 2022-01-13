package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/types"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

const (
	BackgroundCheckStatusQuery = "background-check-status"
)

type backgroundCheckWorkflow struct {
	types.BackgroundCheckState
	checkID      string
	checkFutures map[string]workflow.Future
	logger       log.Logger
}

func newBackgroundCheckWorkflow(ctx workflow.Context, state *types.BackgroundCheckState) (*backgroundCheckWorkflow, error) {
	w := backgroundCheckWorkflow{
		BackgroundCheckState: *state,
		checkID:              workflow.GetInfo(ctx).WorkflowExecution.RunID,
		checkFutures:         make(map[string]workflow.Future),
		logger:               workflow.GetLogger(ctx),
	}

	err := workflow.SetQueryHandler(ctx, BackgroundCheckStatusQuery, func() (types.BackgroundCheckState, error) {
		return w.BackgroundCheckState, nil
	})
	return &w, err
}

func (w *backgroundCheckWorkflow) pushStatus(ctx workflow.Context, status string) error {
	return workflow.UpsertSearchAttributes(
		ctx,
		map[string]interface{}{
			"BackgroundCheckStatus": status,
		},
	)
}

func (w *backgroundCheckWorkflow) waitForAccept(ctx workflow.Context, email string) (*types.AcceptSubmission, error) {
	var r types.AcceptSubmission

	err := w.pushStatus(ctx, "pending_accept")
	if err != nil {
		return &r, err
	}

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: AcceptWorkflowID(email),
	})
	consentWF := workflow.ExecuteChildWorkflow(ctx, Accept, types.AcceptWorkflowInput{
		Email: email,
	})
	err = consentWF.Get(ctx, &r)

	return &r, err
}

func (w *backgroundCheckWorkflow) ssnTrace(ctx workflow.Context) (*types.SSNTraceWorkflowResult, error) {
	var r types.SSNTraceWorkflowResult

	ssnTrace := workflow.ExecuteChildWorkflow(
		ctx,
		SSNTrace,
		types.SSNTraceWorkflowInput{FullName: w.CandidateDetails.FullName, SSN: w.CandidateDetails.SSN},
	)

	err := ssnTrace.Get(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &r, err
}

func (w *backgroundCheckWorkflow) sendDeclineEmail(ctx workflow.Context, email string) error {
	w.pushStatus(ctx, "declined")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendDeclineEmail, types.SendDeclineEmailInput{Email: email})
	return f.Get(ctx, nil)
}

func (w *backgroundCheckWorkflow) sendReportEmail(ctx workflow.Context, email string) error {
	w.pushStatus(ctx, "completed")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendReportEmail, types.SendReportEmailInput{Email: email, State: w.BackgroundCheckState})
	return f.Get(ctx, nil)
}

func (w *backgroundCheckWorkflow) startCheck(ctx workflow.Context, name string, checkWorkflow interface{}, checkInputs ...interface{}) {
	f := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: CheckWorkflowID(w.Email, name),
		}),
		checkWorkflow,
		checkInputs...,
	)
	w.checkFutures[name] = f
}

func (w *backgroundCheckWorkflow) waitForChecks(ctx workflow.Context) {
	for name, f := range w.checkFutures {
		var r interface{}

		err := f.Get(ctx, &r)
		if err != nil {
			w.logger.Error("Search failed", "name", name, "error", err)
			continue
		}

		w.Checks[name] = r
	}
}

// @@@SNIPSTART background-checks-main-workflow-definition
func BackgroundCheck(ctx workflow.Context, input *types.BackgroundCheckWorkflowInput) error {
	w, err := newBackgroundCheckWorkflow(
		ctx,
		&types.BackgroundCheckState{
			Email:  input.Email,
			Tier:   input.Package,
			Checks: make(map[string]interface{}),
		},
	)
	if err != nil {
		return err
	}

	response, err := w.waitForAccept(ctx, w.Email)
	if err != nil {
		return err
	}

	w.Accepted = response.Accepted

	if !w.Accepted {
		return w.sendDeclineEmail(ctx, activities.HiringManagerEmail)
	}

	w.CandidateDetails = response.CandidateDetails

	err = w.pushStatus(ctx, "running")
	if err != nil {
		return err
	}

	t, err := w.ssnTrace(ctx)
	if err != nil {
		return err
	}

	w.SSNTrace = t

	if !w.SSNTrace.SSNIsValid {
		return w.sendReportEmail(ctx, activities.HiringManagerEmail)
	}

	w.startCheck(
		ctx,
		"FederalCriminalSearch",
		FederalCriminalSearch,
		types.FederalCriminalSearchWorkflowInput{FullName: w.CandidateDetails.FullName, Address: w.CandidateDetails.Address},
	)

	if w.Tier == "full" {
		w.startCheck(
			ctx,
			"StateCriminalSearch",
			StateCriminalSearch,
			types.StateCriminalSearchWorkflowInput{FullName: w.CandidateDetails.FullName, SSNTraceResult: w.SSNTrace.KnownAddresses},
		)
		w.startCheck(
			ctx,
			"MotorVehicleIncidentSearch",
			MotorVehicleIncidentSearch,
			types.MotorVehicleIncidentSearchWorkflowInput{FullName: w.CandidateDetails.FullName, Address: w.CandidateDetails.Address},
		)
		if w.CandidateDetails.Employer != "" {
			w.startCheck(
				ctx,
				"EmploymentVerification",
				EmploymentVerification,
				types.EmploymentVerificationWorkflowInput{CandidateDetails: w.CandidateDetails},
			)
		}
	}

	w.waitForChecks(ctx)

	return w.sendReportEmail(ctx, activities.HiringManagerEmail)
}

// @@@SNIPEND
