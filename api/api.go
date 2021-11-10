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

	"github.com/gorilla/mux"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/client"
)

const DefaultEndpoint = "localhost:8081"
const TaskQueue = "background-checks-main"

type CheckListItem struct {
	Email string
	Tier  string
}

type CheckCreateInput struct {
	Email string
	Tier  string
}

type CheckSaveSearchResultInput struct {
	Type                             string
	FederalCriminalSearchResult      types.FederalCriminalSearchResult
	StateCriminalSearchResult        types.StateCriminalSearchResult
	MotorVehicleIncidentSearchResult types.MotorVehicleIncidentSearchResult
}

func handleCheckList(w http.ResponseWriter, r *http.Request) {
	checks := []CheckListItem{}

	// client.ListOpenWorkflowExecutions

	json.NewEncoder(w).Encode(checks)
}

func handleCheckCreate(w http.ResponseWriter, r *http.Request) {
	var input CheckCreateInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	_, err = c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
		},
		workflows.BackgroundCheck,
		types.BackgroundCheckInput{
			Email: input.Email,
			Tier:  input.Tier,
		},
	)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleCheckDetails(w http.ResponseWriter, r *http.Request) {
}

func handleCheckReport(w http.ResponseWriter, r *http.Request) {
}

func handleCheckConsent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	token, err := base64.StdEncoding.DecodeString(vars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result types.ConsentResult

	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	err = c.CompleteActivity(context.Background(), token, result, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCheckDecline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	token, err := base64.StdEncoding.DecodeString(vars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	// TODO: This needs to be a defined error type so that we don't retry.
	err = c.CompleteActivity(context.Background(), token, nil, fmt.Errorf("declined"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCheckCancel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	err = c.CancelWorkflow(context.Background(), id, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCandidateTodoList(w http.ResponseWriter, r *http.Request) {
}

func handleResearcherTodoList(w http.ResponseWriter, r *http.Request) {
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

	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	err = c.CompleteActivity(context.Background(), token, result.Result(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/checks", handleCheckList).Name("checks_list")
	r.HandleFunc("/checks", handleCheckCreate).Methods("POST").Name("checks_create")
	r.HandleFunc("/checks/{id}", handleCheckDetails).Name("check")
	r.HandleFunc("/checks/{id}/cancel", handleCheckCancel).Methods("POST").Name("check_cancel")
	r.HandleFunc("/checks/{id}/report", handleCheckReport).Name("check_report")
	r.HandleFunc("/checks/{token}/consent", handleCheckConsent).Methods("POST").Name("check_consent")
	r.HandleFunc("/checks/{token}/decline", handleCheckDecline).Methods("POST").Name("check_decline")
	r.HandleFunc("/checks/{token}/search", handleSaveSearchResult).Methods("POST").Name("research_save")
	r.HandleFunc("/todos/candidate/{email}", handleCandidateTodoList).Name("todos_candidate")
	r.HandleFunc("/todos/researcher/{email}", handleResearcherTodoList).Name("todos_researcher")

	return r
}

func Run() {
	srv := &http.Server{
		Handler: Router(),
		Addr:    DefaultEndpoint,
	}

	log.Fatal(srv.ListenAndServe())
}
