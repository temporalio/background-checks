package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"go.temporal.io/api/enums/v1"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"

	"github.com/temporalio/background-checks/workflows"
)

const (
	TaskQueue       = "background-checks-main"
	DefaultEndpoint = "0.0.0.0:8081"
)

type handlers struct {
	temporalClient client.Client
}

func getBackgroundCheckCandidateEmail(we *workflowpb.WorkflowExecutionInfo) (string, error) {
	var email string

	attrs := we.GetSearchAttributes().GetIndexedFields()

	err := converter.GetDefaultDataConverter().FromPayload(attrs["CandidateEmail"], &email)

	return email, err
}

func getBackgroundCheckStatus(we *workflowpb.WorkflowExecutionInfo) (string, error) {
	var status string

	attrs := we.GetSearchAttributes().GetIndexedFields()

	err := converter.GetDefaultDataConverter().FromPayload(attrs["BackgroundCheckStatus"], &status)

	return status, err
}

func presentBackgroundCheck(we *workflowpb.WorkflowExecutionInfo) (BackgroundCheck, error) {
	var result BackgroundCheck

	result.ID = we.Execution.RunId

	email, err := getBackgroundCheckCandidateEmail(we)
	if err != nil {
		return result, err
	}
	result.Email = email

	checkStatus, err := getBackgroundCheckStatus(we)
	if err != nil {
		return result, err
	}

	switch we.Status {
	case enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
		result.Status = checkStatus
	case enums.WORKFLOW_EXECUTION_STATUS_TIMED_OUT,
		enums.WORKFLOW_EXECUTION_STATUS_FAILED:
		result.Status = "failed"
	case enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		if checkStatus == "declined" {
			result.Status = "declined"
		} else {
			result.Status = "completed"
		}
	case enums.WORKFLOW_EXECUTION_STATUS_TERMINATED:
		result.Status = "terminated"
	case enums.WORKFLOW_EXECUTION_STATUS_CANCELED:
		result.Status = "cancelled"
	default:
		result.Status = "unknown"
	}

	return result, nil
}

type listWorkflowFilters struct {
	Email  string
	Status string
}

func queryForFilters(filters listWorkflowFilters) (string, error) {
	q := []string{
		"WorkflowType = 'BackgroundCheck'",
	}

	if filters.Email != "" {
		q = append(q, candidateQuery(filters.Email))
	}
	if filters.Status != "" {
		f, err := statusQuery(filters.Status)
		if err != nil {
			return "", err
		}
		q = append(q, f)
	}

	query := strings.Join(q, " AND ")

	return query, nil
}

func candidateQuery(email string) string {
	return fmt.Sprintf("CandidateEmail = '%s'", email)
}

func statusQuery(status string) (string, error) {
	switch status {
	case "pending_accept":
		return fmt.Sprintf("ExecutionStatus = 'Running' AND BackgroundCheckStatus = '%s'", status), nil
	case "running":
		return fmt.Sprintf("ExecutionStatus = 'Running' AND BackgroundCheckStatus = '%s'", status), nil
	case "completed":
		return fmt.Sprintf("ExecutionStatus = 'Completed' AND BackgroundCheckStatus = '%s'", status), nil
	case "declined":
		return fmt.Sprintf("ExecutionStatus = 'Completed' AND BackgroundCheckStatus = '%s'", status), nil
	case "failed":
		return "ExecutionStatus = 'Failed'", nil
	case "terminated":
		return "ExecutionStatus = 'Terminated'", nil
	case "cancelled":
		return "ExecutionStatus = 'Cancelled'", nil
	default:
		return "", fmt.Errorf("unknown status: %s", status)
	}
}

func (h *handlers) listWorkflows(ctx context.Context, filters listWorkflowFilters) ([]*workflowpb.WorkflowExecutionInfo, error) {
	var executions []*workflowpb.WorkflowExecutionInfo
	var nextPageToken []byte

	query, err := queryForFilters(filters)
	if err != nil {
		return executions, err
	}

	for {
		resp, err := h.temporalClient.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			PageSize:      10,
			NextPageToken: nextPageToken,
			Query:         query,
		})
		if err != nil {
			return executions, err
		}

		executions = append(executions, resp.Executions...)
		if len(resp.NextPageToken) == 0 {
			return executions, nil
		}
		nextPageToken = resp.NextPageToken
	}
}

func (h *handlers) handleCheckList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters := listWorkflowFilters{
		Email:  query.Get("email"),
		Status: query.Get("status"),
	}

	wfs, err := h.listWorkflows(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checks := make([]BackgroundCheck, len(wfs))
	for i, wf := range wfs {
		check, err := presentBackgroundCheck(wf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		checks[i] = check
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checks)
}

func (h *handlers) handleCheckCreate(w http.ResponseWriter, r *http.Request) {
	var input workflows.BackgroundCheckWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.temporalClient.ExecuteWorkflow(
		r.Context(),
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
			ID:        workflows.BackgroundCheckWorkflowID(input.Email),
			SearchAttributes: map[string]interface{}{
				"CandidateEmail": input.Email,
			},
		},
		workflows.BackgroundCheck,
		&input,
	)

	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handlers) handleCheckStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]

	v, err := h.temporalClient.QueryWorkflow(
		r.Context(),
		workflows.BackgroundCheckWorkflowID(email),
		"",
		workflows.BackgroundCheckStatusQuery,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result workflows.BackgroundCheckState
	err = v.Get(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *handlers) handleCheckReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	wfid, runid, err := workflows.WorkflowFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enc, err := h.temporalClient.QueryWorkflow(
		r.Context(),
		wfid,
		runid,
		workflows.BackgroundCheckStatusQuery,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result workflows.BackgroundCheckState
	err = enc.Get(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *handlers) handleAccept(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	wfid, runid, err := workflows.WorkflowFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result workflows.AcceptSubmissionSignal

	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result.Accepted = true

	err = h.temporalClient.SignalWorkflow(
		r.Context(),
		wfid,
		runid,
		workflows.AcceptSubmissionSignalName,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handlers) handleDecline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	wfid, runid, err := workflows.WorkflowFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := workflows.AcceptSubmissionSignal{
		Accepted: false,
	}

	err = h.temporalClient.SignalWorkflow(
		r.Context(),
		wfid,
		runid,
		workflows.AcceptSubmissionSignalName,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handlers) handleEmploymentVerificationDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	wfid, runid, err := workflows.WorkflowFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enc, err := h.temporalClient.QueryWorkflow(
		r.Context(),
		wfid,
		runid,
		workflows.EmploymentVerificationDetailsQuery,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result workflows.CandidateDetails
	err = enc.Get(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *handlers) handleEmploymentVerificationSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	wfid, runid, err := workflows.WorkflowFromToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input workflows.EmploymentVerificationSubmissionSignal

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := input

	err = h.temporalClient.SignalWorkflow(
		r.Context(),
		wfid,
		runid,
		workflows.EmploymentVerificationSubmissionSignalName,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *handlers) handleCheckCancel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]
	wfid := workflows.BackgroundCheckWorkflowID(email)
	id := vars["id"]

	err := h.temporalClient.CancelWorkflow(r.Context(), wfid, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()

	h := handlers{temporalClient: c}

	r.HandleFunc("/checks", h.handleCheckList).Methods("GET").Name("checks_list")
	r.HandleFunc("/checks", h.handleCheckCreate).Methods("POST").Name("checks_create")
	r.HandleFunc("/checks/{email}/{id}/cancel", h.handleCheckCancel).Methods("POST").Name("check_cancel")
	r.HandleFunc("/checks/{email}", h.handleCheckStatus).Methods("GET").Name("check")

	r.HandleFunc("/checks/{token}/accept", h.handleAccept).Methods("POST").Name("accept")
	r.HandleFunc("/checks/{token}/decline", h.handleDecline).Methods("POST").Name("decline")

	r.HandleFunc("/checks/{token}/employment", h.handleEmploymentVerificationDetails).Methods("GET").Name("employmentverify_details")
	r.HandleFunc("/checks/{token}/employment", h.handleEmploymentVerificationSubmission).Methods("POST").Name("employmentverify")

	r.HandleFunc("/checks/{token}/report", h.handleCheckReport).Methods("GET").Name("check_report")

	return r
}
