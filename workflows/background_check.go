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

// backgroundCheckWorkflow represents the state for a background check workflow execution.
type backgroundCheckWorkflow struct {
	types.BackgroundCheckState
	checkID       string
	searchFutures map[string]workflow.Future
	logger        log.Logger
}

// newBackgroundCheckWorkflow initializes a backgroundCheckWorkflow struct.
func newBackgroundCheckWorkflow(ctx workflow.Context, state *types.BackgroundCheckState) *backgroundCheckWorkflow {
	return &backgroundCheckWorkflow{
		BackgroundCheckState: *state,
		checkID:              workflow.GetInfo(ctx).WorkflowExecution.RunID,
		searchFutures:        make(map[string]workflow.Future),
		logger:               workflow.GetLogger(ctx),
	}
}

// pushStatus updates the BackgroundCheckStatus search attribute for a background check workflow execution.
func (w *backgroundCheckWorkflow) pushStatus(ctx workflow.Context, status string) error {
	return workflow.UpsertSearchAttributes(
		ctx,
		map[string]interface{}{
			"BackgroundCheckStatus": status,
		},
	)
}

// waitForAccept waits for the candidate to accept or decline the background check.
// If the candidate accepted, the response will include their personal information.
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

// ssnTrace runs an SSN trace.
// This will tell us if the SSN the candidate gave us is valid.
// It also provides us with a list of addresses that the candidate is linked to in the SSN system.
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

// sendDeclineEmail sends an email to the Hiring Manager informing them the candidate declined the background check.
func (w *backgroundCheckWorkflow) sendDeclineEmail(ctx workflow.Context, email string) error {
	w.pushStatus(ctx, "declined")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendDeclineEmail, types.SendDeclineEmailInput{Email: w.Email})
	return f.Get(ctx, nil)
}

// sendReportEmail sends an email to the Hiring Manager with a link to the report page for the background check.
func (w *backgroundCheckWorkflow) sendReportEmail(ctx workflow.Context, email string) error {
	w.pushStatus(ctx, "completed")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendReportEmail, types.SendReportEmailInput{Email: w.Email, Token: TokenForWorkflow(ctx)})
	return f.Get(ctx, nil)
}

// startSearch starts a child workflow to perform one of the searches that make up the background check.
func (w *backgroundCheckWorkflow) startSearch(ctx workflow.Context, name string, searchWorkflow interface{}, searchInputs ...interface{}) {
	f := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: SearchWorkflowID(w.Email, name),
		}),
		searchWorkflow,
		searchInputs...,
	)
	// Record the future for the search so we can collect the results later
	w.searchFutures[name] = f
}

// waitForSearches waits for all of our searches to complete and collects the results.
func (w *backgroundCheckWorkflow) waitForSearches(ctx workflow.Context) {
	for name, f := range w.searchFutures {
		var r interface{}

		err := f.Get(ctx, &r)
		if err != nil {
			w.logger.Error("Search failed", "name", name, "error", err)
			// Record an error for the search so we can include it in the report.
			w.SearchErrors[name] = err.Error()
			continue
		}
		// Record the result of the search so we can use it in the report.
		w.SearchResults[name] = r
	}
}

// @@@SNIPSTART background-checks-main-workflow-definition

// BackgroundCheck is a Workflow Definition that calls for the execution of a variable set of Activities and Child Workflows.
// This is the main entry point of the application.
// It accepts an email address as the input.
// All other personal information for the Candidate is provided when they accept the Background Check.
func BackgroundCheck(ctx workflow.Context, input *types.BackgroundCheckWorkflowInput) (*types.BackgroundCheckWorkflowResult, error) {
	w := newBackgroundCheckWorkflow(
		ctx,
		&types.BackgroundCheckState{
			Email:         input.Email,
			Tier:          input.Tier,
			SearchResults: make(map[string]interface{}),
			SearchErrors:  make(map[string]string),
		},
	)

	// The query returns the status of a background check and is used by the API to build the report at the end.
	err := workflow.SetQueryHandler(ctx, BackgroundCheckStatusQuery, func() (types.BackgroundCheckState, error) {
		return w.BackgroundCheckState, nil
	})
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// Send the candidate an email asking them to accept or decline the background check.
	response, err := w.waitForAccept(ctx, w.Email)
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	w.Accepted = response.Accepted

	// If the candidate declined the check, let the hiring manager know and then end the workflow.
	if !w.Accepted {
		return &w.BackgroundCheckState, w.sendDeclineEmail(ctx, activities.HiringManagerEmail)
	}

	w.CandidateDetails = response.CandidateDetails

	// Update our status search attribute. This is used by our API to filter the background check list if requested.
	err = w.pushStatus(ctx, "running")
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// Run an SSN trace on the SSN the candidate provided when accepting the background check.
	w.SSNTrace, err = w.ssnTrace(ctx)
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// If the SSN the candidate gave us was not valid then send a report email to the Hiring Manager and end the workflow.
	// In this case all the searches are skipped.
	if !w.SSNTrace.SSNIsValid {
		return &w.BackgroundCheckState, w.sendReportEmail(ctx, activities.HiringManagerEmail)
	}

	// Start the main searches, these are run in parallel as they do not depend on each other.

	var primaryAddress string
	if len(w.SSNTrace.KnownAddresses) > 0 {
		primaryAddress = w.SSNTrace.KnownAddresses[0]
	}

	// We always run the FederalCriminalSearch
	w.startSearch(
		ctx,
		"FederalCriminalSearch",
		FederalCriminalSearch,
		types.FederalCriminalSearchWorkflowInput{FullName: w.CandidateDetails.FullName, KnownAddresses: w.SSNTrace.KnownAddresses},
	)

	// If the background check is on the full tier we run more searches
	if w.Tier == "full" {
		w.startSearch(
			ctx,
			"StateCriminalSearch",
			StateCriminalSearch,
			types.StateCriminalSearchWorkflowInput{FullName: w.CandidateDetails.FullName, KnownAddresses: w.SSNTrace.KnownAddresses},
		)
		w.startSearch(
			ctx,
			"MotorVehicleIncidentSearch",
			MotorVehicleIncidentSearch,
			types.MotorVehicleIncidentSearchWorkflowInput{FullName: w.CandidateDetails.FullName, Address: primaryAddress},
		)

		// Verify their employment if they provided an employer
		if w.CandidateDetails.Employer != "" {
			w.startSearch(
				ctx,
				"EmploymentVerification",
				EmploymentVerification,
				types.EmploymentVerificationWorkflowInput{CandidateDetails: w.CandidateDetails},
			)
		}
	}

	// Wait for all of our searches to complete.
	w.waitForSearches(ctx)

	// Send the report email to the Hiring Manager.
	return &w.BackgroundCheckState, w.sendReportEmail(ctx, activities.HiringManagerEmail)
}

// @@@SNIPEND
