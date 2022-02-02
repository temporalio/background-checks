package workflows

import (
	"time"

	"github.com/temporalio/background-checks/activities"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

const (
	BackgroundCheckStatusQuery = "background-check-status"
)

type BackgroundCheckWorkflowInput struct {
	Email string
	Tier  string
}

type BackgroundCheckState struct {
	Email            string
	Status           string
	Tier             string
	AcceptSubmission *AcceptSubmission
	SSNTrace         *SSNTraceWorkflowResult
	SearchSet        []string
	SearchResults    map[string]interface{}
	SearchErrors     map[string]string
}

type BackgroundCheckWorkflowResult = BackgroundCheckState

// backgroundCheckWorkflow represents the state for a background check workflow execution.
type backgroundCheckWorkflow struct {
	BackgroundCheckState
	checkID       string
	searchFutures map[string]workflow.Future
	logger        log.Logger
}

// newBackgroundCheckWorkflow initializes a backgroundCheckWorkflow struct.
func newBackgroundCheckWorkflow(ctx workflow.Context, state *BackgroundCheckState) *backgroundCheckWorkflow {
	return &backgroundCheckWorkflow{
		BackgroundCheckState: *state,
		checkID:              workflow.GetInfo(ctx).WorkflowExecution.RunID,
		searchFutures:        make(map[string]workflow.Future),
		logger:               workflow.GetLogger(ctx),
	}
}

// pushStatus updates the BackgroundCheckStatus search attribute for a background check workflow execution.
func (w *backgroundCheckWorkflow) pushStatus(ctx workflow.Context, status string) error {
	w.Status = status
	return workflow.UpsertSearchAttributes(
		ctx,
		map[string]interface{}{
			"BackgroundCheckStatus": status,
		},
	)
}

// waitForAccept waits for the candidate to accept or decline the background check.
// If the candidate accepted, the response will include their personal information.
func (w *backgroundCheckWorkflow) waitForAccept(ctx workflow.Context, email string) (*AcceptSubmission, error) {
	var r AcceptSubmission

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: AcceptWorkflowID(email),
	})
	consentWF := workflow.ExecuteChildWorkflow(ctx, Accept, AcceptWorkflowInput{
		Email: email,
	})
	err := consentWF.Get(ctx, &r)

	return &r, err
}

// ssnTrace runs an SSN trace.
// This will tell us if the SSN the candidate gave us is valid.
// It also provides us with a list of addresses that the candidate is linked to in the SSN system.
func (w *backgroundCheckWorkflow) ssnTrace(ctx workflow.Context, name string, ssn string) (*SSNTraceWorkflowResult, error) {
	var r SSNTraceWorkflowResult

	ssnTrace := workflow.ExecuteChildWorkflow(
		ctx,
		SSNTrace,
		SSNTraceWorkflowInput{FullName: name, SSN: ssn},
	)

	err := ssnTrace.Get(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &r, err
}

// sendDeclineEmail sends an email to the Hiring Manager informing them the candidate declined the background check.
func (w *backgroundCheckWorkflow) sendDeclineEmail(ctx workflow.Context, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendDeclineEmail, activities.SendDeclineEmailInput{Email: w.Email})
	return f.Get(ctx, nil)
}

// sendReportEmail sends an email to the Hiring Manager with a link to the report page for the background check.
func (w *backgroundCheckWorkflow) sendReportEmail(ctx workflow.Context, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	f := workflow.ExecuteActivity(ctx, a.SendReportEmail, activities.SendReportEmailInput{Email: w.Email, Token: TokenForWorkflow(ctx)})
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
	// Add the name of the search to the search set. This is used by the UI so it can tell what is outstanding.
	w.SearchSet = append(w.SearchSet, name)
	// Record the future for the search so we can collect the results later
	w.searchFutures[name] = f
}

// waitForSearches waits for all of our searches to complete and collects the results.
func (w *backgroundCheckWorkflow) waitForSearches(ctx workflow.Context) {
	s := workflow.NewSelector(ctx)
	for name, f := range w.searchFutures {
		name := name

		s.AddFuture(f, func(rf workflow.Future) {
			var r interface{}

			err := rf.Get(ctx, &r)

			if err != nil {
				w.logger.Error("Search failed", "name", name, "error", err)
				// Record an error for the search so we can include it in the report.
				w.SearchErrors[name] = err.Error()
				return
			}

			// Record the result of the search so we can use it in the report.
			w.SearchResults[name] = r
		})
	}

	for range w.searchFutures {
		s.Select(ctx)
	}
}

// @@@SNIPSTART background-checks-main-workflow-definition

// BackgroundCheck is a Workflow Definition that calls for the execution of a variable set of Activities and Child Workflows.
// This is the main entry point of the application.
// It accepts an email address as the input.
// All other personal information for the Candidate is provided when they accept the Background Check.
func BackgroundCheck(ctx workflow.Context, input *BackgroundCheckWorkflowInput) (*BackgroundCheckWorkflowResult, error) {
	w := newBackgroundCheckWorkflow(
		ctx,
		&BackgroundCheckState{
			Email:         input.Email,
			Tier:          input.Tier,
			SearchResults: make(map[string]interface{}),
			SearchErrors:  make(map[string]string),
		},
	)

	// The query returns the status of a background check and is used by the API to build the report at the end.
	err := workflow.SetQueryHandler(ctx, BackgroundCheckStatusQuery, func() (BackgroundCheckState, error) {
		return w.BackgroundCheckState, nil
	})
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	err = w.pushStatus(ctx, "pending_accept")
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// Send the candidate an email asking them to accept or decline the background check.
	w.AcceptSubmission, err = w.waitForAccept(ctx, w.Email)
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// If the candidate declined the check, let the hiring manager know and then end the workflow.
	if !w.AcceptSubmission.Accepted {
		err = w.pushStatus(ctx, "declined")
		if err != nil {
			return &w.BackgroundCheckState, err
		}

		return &w.BackgroundCheckState, w.sendDeclineEmail(ctx, activities.HiringManagerEmail)
	}

	candidateDetails := w.AcceptSubmission.CandidateDetails

	// Update our status search attribute. This is used by our API to filter the background check list if requested.
	err = w.pushStatus(ctx, "running")
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// Run an SSN trace on the SSN the candidate provided when accepting the background check.
	w.SSNTrace, err = w.ssnTrace(ctx, candidateDetails.FullName, candidateDetails.SSN)
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// If the SSN the candidate gave us was not valid then send a report email to the Hiring Manager and end the workflow.
	// In this case all the searches are skipped.
	if !w.SSNTrace.SSNIsValid {
		err = w.pushStatus(ctx, "completed")
		if err != nil {
			return &w.BackgroundCheckState, err
		}

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
		FederalCriminalSearchWorkflowInput{FullName: candidateDetails.FullName, KnownAddresses: w.SSNTrace.KnownAddresses},
	)

	// If the background check is on the full tier we run more searches
	if w.Tier == "full" {
		w.startSearch(
			ctx,
			"StateCriminalSearch",
			StateCriminalSearch,
			StateCriminalSearchWorkflowInput{FullName: candidateDetails.FullName, KnownAddresses: w.SSNTrace.KnownAddresses},
		)
		w.startSearch(
			ctx,
			"MotorVehicleIncidentSearch",
			MotorVehicleIncidentSearch,
			MotorVehicleIncidentSearchWorkflowInput{FullName: candidateDetails.FullName, Address: primaryAddress},
		)

		// Verify their employment if they provided an employer
		if candidateDetails.Employer != "" {
			w.startSearch(
				ctx,
				"EmploymentVerification",
				EmploymentVerification,
				EmploymentVerificationWorkflowInput{CandidateDetails: candidateDetails},
			)
		}
	}

	// Wait for all of our searches to complete.
	w.waitForSearches(ctx)

	err = w.pushStatus(ctx, "completed")
	if err != nil {
		return &w.BackgroundCheckState, err
	}

	// Send the report email to the Hiring Manager.
	return &w.BackgroundCheckState, w.sendReportEmail(ctx, activities.HiringManagerEmail)
}

// @@@SNIPEND
