/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	sdkclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/queries"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
)

const DefaultEndpoint = "localhost:8081"

var client sdkclient.Client

func getClient() (sdkclient.Client, error) {
	if client != nil {
		return client, nil
	}

	c, err := sdkclient.NewClient(sdkclient.Options{})
	if err != nil {
		return nil, err
	}
	client = c

	return c, nil
}

func executeWorkflow(options sdkclient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (sdkclient.WorkflowRun, error) {
	c, err := getClient()
	if err != nil {
		return nil, err
	}

	options.TaskQueue = config.TaskQueue

	return c.ExecuteWorkflow(
		context.Background(),
		options,
		workflows.BackgroundCheck,
		args...,
	)
}

func cancelWorkflow(wid string) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.CancelWorkflow(context.Background(), wid, "")
}

func completeActivity(token []byte, result interface{}, activityErr error) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.CompleteActivity(context.Background(), token, result, activityErr)
}

func queryWorkflow(wid string, queryType string, args ...interface{}) (converter.EncodedValue, error) {
	c, err := getClient()
	if err != nil {
		return nil, err
	}

	return c.QueryWorkflow(
		context.Background(),
		wid,
		"",
		queryType,
		args...,
	)
}

func signalWorkflow(wid string, signalName string, signalArg interface{}) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.SignalWorkflow(context.Background(), wid, "", signalName, signalArg)
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
	switch types.BackgroundCheckStatusFromString(status) {
	case types.BackgroundCheckStatusPendingAccept:
		return fmt.Sprintf("ExecutionStatus = 'Running' AND BackgroundCheckStatus = '%s'", status), nil
	case types.BackgroundCheckStatusRunning:
		return fmt.Sprintf("ExecutionStatus = 'Running' AND BackgroundCheckStatus = '%s'", status), nil
	case types.BackgroundCheckStatusCompleted:
		return fmt.Sprintf("ExecutionStatus = 'Completed' AND BackgroundCheckStatus = '%s'", status), nil
	case types.BackgroundCheckStatusDeclined:
		return fmt.Sprintf("ExecutionStatus = 'Completed' AND BackgroundCheckStatus = '%s'", status), nil
	case types.BackgroundCheckStatusFailed:
		return "ExecutionStatus = 'Failed'", nil
	case types.BackgroundCheckStatusTerminated:
		return "ExecutionStatus = 'Terminated'", nil
	case types.BackgroundCheckStatusCancelled:
		return "ExecutionStatus = 'Cancelled'", nil
	default:
		return "", fmt.Errorf("unknown status: %s", status)
	}
}

func listWorkflows(filters listWorkflowFilters) ([]*workflowpb.WorkflowExecutionInfo, error) {
	var executions []*workflowpb.WorkflowExecutionInfo
	var nextPageToken []byte

	c, err := getClient()
	if err != nil {
		return executions, err
	}

	ctx := context.Background()

	query, err := queryForFilters(filters)
	if err != nil {
		return executions, err
	}

	for hasMore := true; hasMore; hasMore = len(nextPageToken) > 0 {
		resp, err := c.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
			PageSize:      10,
			NextPageToken: nextPageToken,
			Query:         query,
		})
		if err != nil {
			return executions, err
		}

		executions = append(executions, resp.Executions...)
		nextPageToken = resp.NextPageToken
	}

	return executions, nil
}

func handleCheckList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters := listWorkflowFilters{
		Email:  query.Get("email"),
		Status: query.Get("status"),
	}

	wfs, err := listWorkflows(filters)
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

	json.NewEncoder(w).Encode(checks)
}

func handleCheckCreate(w http.ResponseWriter, r *http.Request) {
	var input types.BackgroundCheckWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = executeWorkflow(
		sdkclient.StartWorkflowOptions{
			ID: mappings.BackgroundCheckWorkflowID(input.Email),
			SearchAttributes: map[string]interface{}{
				"CandidateEmail": input.Email,
			},
		},
		workflows.BackgroundCheck,
		input,
	)

	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleCheckStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]

	v, err := queryWorkflow(
		mappings.BackgroundCheckWorkflowID(email),
		queries.BackgroundCheckStatus,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result types.BackgroundCheckState
	err = v.Get(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleCheckReport(w http.ResponseWriter, r *http.Request) {
}

func handleAccept(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var result types.AcceptSubmissionSignal
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = signalWorkflow(
		mappings.AcceptWorkflowID(id),
		signals.AcceptSubmission,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDecline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	result := types.AcceptSubmissionSignal{
		Accept: types.Accept{Accept: false},
	}

	err := signalWorkflow(
		mappings.AcceptWorkflowID(id),
		signals.AcceptSubmission,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCheckCancel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := cancelWorkflow(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCandidateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]

	filters := listWorkflowFilters{
		Email:  email,
		Status: "pending_consent",
	}

	wfs, err := listWorkflows(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(wfs) == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	check, err := presentBackgroundCheck(wfs[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(check)
}

func handleResearcherStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]

	v, err := queryWorkflow(
		mappings.ResearcherWorkflowID(email),
		queries.ResearcherTodosList,
	)
	if _, ok := err.(*serviceerror.NotFound); ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []types.ResearcherTodo
	err = v.Get(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleSaveSearchResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	token, err := base64.StdEncoding.DecodeString(vars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result types.SearchResult
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = completeActivity(token, result.Result(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/checks", handleCheckList).Methods("GET").Name("checks_list")
	r.HandleFunc("/checks", handleCheckCreate).Methods("POST").Name("checks_create")
	r.HandleFunc("/checks/{email}", handleCheckStatus).Name("check")
	r.HandleFunc("/checks/{email}/cancel", handleCheckCancel).Methods("POST").Name("check_cancel")
	r.HandleFunc("/checks/{email}/report", handleCheckReport).Name("check_report")
	r.HandleFunc("/checks/{id}/accept", handleAccept).Methods("POST").Name("accept")
	r.HandleFunc("/checks/{id}/decline", handleDecline).Methods("POST").Name("decline")
	r.HandleFunc("/checks/{token}/search", handleSaveSearchResult).Methods("POST").Name("research_save")
	r.HandleFunc("/candidate/{email}", handleCandidateStatus).Name("candidate")
	r.HandleFunc("/research/{email}", handleResearcherStatus).Name("research")

	return r
}

func Run() {
	srv := &http.Server{
		Handler: Router(),
		Addr:    DefaultEndpoint,
	}

	log.Fatal(srv.ListenAndServe())
}
