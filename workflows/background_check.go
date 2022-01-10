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
	checkID       string
	ctx           workflow.Context
	checkSelector workflow.Selector
	logger        log.Logger
}

func newBackgroundCheckWorkflow(ctx workflow.Context, state *types.BackgroundCheckState) (*backgroundCheckWorkflow, error) {
	w := backgroundCheckWorkflow{
		BackgroundCheckState: *state,
		checkID:              workflow.GetInfo(ctx).WorkflowExecution.RunID,
		ctx:                  ctx,
		checkSelector:        workflow.NewSelector(ctx),
		logger:               workflow.GetLogger(ctx),
	}

	w.Checks = make(map[string]interface{})

	err := workflow.SetQueryHandler(ctx, BackgroundCheckStatusQuery, func() (types.BackgroundCheckState, error) {
		return w.BackgroundCheckState, nil
	})
	return &w, err
}

func (w *backgroundCheckWorkflow) pushStatus(status string) error {
	return workflow.UpsertSearchAttributes(
		w.ctx,
		map[string]interface{}{
			"BackgroundCheckStatus": status,
		},
	)
}

func (w *backgroundCheckWorkflow) waitForAccept(email string) (*types.AcceptSubmission, error) {
	var r types.AcceptSubmission

	err := w.pushStatus("pending_accept")
	if err != nil {
		return &r, err
	}

	ctx := workflow.WithChildOptions(w.ctx, workflow.ChildWorkflowOptions{
		WorkflowID: AcceptWorkflowID(email),
	})
	consentWF := workflow.ExecuteChildWorkflow(ctx, Accept, types.AcceptWorkflowInput{
		Email: email,
	})
	err = consentWF.Get(ctx, &r)

	return &r, err
}

func (w *backgroundCheckWorkflow) sendDeclineEmail(email string) error {
	w.pushStatus("declined")

	ctx := workflow.WithActivityOptions(w.ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendDeclineEmail, types.SendDeclineEmailInput{Email: email})
	return f.Get(w.ctx, nil)
}

func (w *backgroundCheckWorkflow) sendReportEmail(email string) error {
	w.pushStatus("completed")

	ctx := workflow.WithActivityOptions(w.ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendReportEmail, types.SendReportEmailInput{Email: email, State: w.BackgroundCheckState})
	return f.Get(ctx, nil)
}

func (w *backgroundCheckWorkflow) startCheck(name string, checkWorkflow interface{}, checkInputs ...interface{}) {
	f := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(w.ctx, workflow.ChildWorkflowOptions{
			WorkflowID: CheckWorkflowID(w.Email, name),
		}),
		checkWorkflow,
		checkInputs...,
	)
	w.checkSelector.AddFuture(f, func(f workflow.Future) {
		var result interface{}

		err := f.Get(w.ctx, &result)
		if err != nil {
			w.logger.Error("Search failed", "name", name, "error", err)
		}

		w.Checks[name] = result
	})
}

func (w *backgroundCheckWorkflow) waitForChecks() {
	for w.checkSelector.HasPending() {
		w.checkSelector.Select(w.ctx)
	}
}

// @@@SNIPSTART background-checks-main-workflow-definition
func BackgroundCheck(ctx workflow.Context, input *types.BackgroundCheckWorkflowInput) error {
	w, err := newBackgroundCheckWorkflow(
		ctx,
		&types.BackgroundCheckState{
			Email: input.Email,
			Tier:  input.Package,
		},
	)
	if err != nil {
		return err
	}

	response, err := w.waitForAccept(w.Email)
	if err != nil {
		return err
	}

	w.Accepted = response.Accepted

	if !w.Accepted {
		return w.sendDeclineEmail(activities.HiringManagerEmail)
	}

	w.CandidateDetails = response.CandidateDetails

	err = w.pushStatus("running")
	if err != nil {
		return err
	}

	ssnTrace := workflow.ExecuteChildWorkflow(
		ctx,
		SSNTrace,
		types.SSNTraceWorkflowInput{FullName: w.CandidateDetails.FullName, SSN: w.CandidateDetails.SSN},
	)

	err = ssnTrace.Get(ctx, &w.SSNTrace)
	if err != nil {
		return err
	}

	if !w.SSNTrace.SSNIsValid {
		return w.sendReportEmail(activities.HiringManagerEmail)
	}

	w.startCheck(
		"FederalCriminalSearch",
		FederalCriminalSearch,
		types.FederalCriminalSearchWorkflowInput{FullName: w.CandidateDetails.FullName, Address: w.CandidateDetails.Address},
	)

	if w.Tier == "full" {
		w.startCheck(
			"StateCriminalSearch",
			StateCriminalSearch,
			types.StateCriminalSearchWorkflowInput{FullName: w.CandidateDetails.FullName, SSNTraceResult: w.SSNTrace.KnownAddresses},
		)
		w.startCheck(
			"MotorVehicleIncidentSearch",
			MotorVehicleIncidentSearch,
			types.MotorVehicleIncidentSearchWorkflowInput{FullName: w.CandidateDetails.FullName, Address: w.CandidateDetails.Address},
		)
		w.startCheck(
			"EmploymentVerification",
			EmploymentVerification,
			types.EmploymentVerificationWorkflowInput{CandidateDetails: w.CandidateDetails},
		)
	}

	w.waitForChecks()

	return w.sendReportEmail(activities.HiringManagerEmail)
}

// @@@SNIPEND
