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
	"github.com/temporalio/background-checks/internal"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/client"
)

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
	FederalCriminalSearchResult      internal.FederalCriminalSearchResult
	StateCriminalSearchResult        internal.StateCriminalSearchResult
	MotorVehicleIncidentSearchResult internal.MotorVehicleIncidentSearchResult
}

func handleCheckList(w http.ResponseWriter, r *http.Request) {
	checks := []CheckListItem{}

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
		internal.BackgroundCheckInput{
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

func handleCheckAccept(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	token, err := base64.StdEncoding.DecodeString(vars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result internal.AcceptCheckResult

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

func handleCandidateTasksList(w http.ResponseWriter, r *http.Request) {
}

func handleResearcherTasksList(w http.ResponseWriter, r *http.Request) {
}

func handleSaveSearchResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	token, err := base64.StdEncoding.DecodeString(vars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result CheckSaveSearchResultInput

	if err != nil {
		err = json.NewDecoder(r.Body).Decode(&result)
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

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/checks", handleCheckList)
	r.HandleFunc("/checks", handleCheckCreate).Methods("POST")
	r.HandleFunc("/checks/{id}", handleCheckDetails)
	r.HandleFunc("/checks/{id}/cancel", handleCheckCancel).Methods("POST")
	r.HandleFunc("/checks/{id}/report", handleCheckReport)
	r.HandleFunc("/checks/{token}/accept", handleCheckAccept).Methods("POST")
	r.HandleFunc("/checks/{token}/decline", handleCheckDecline).Methods("POST")
	r.HandleFunc("/checks/{token}/search", handleSaveSearchResult).Methods("POST")
	r.HandleFunc("/tasks/candidate/{email}", handleCandidateTasksList)
	r.HandleFunc("/tasks/researcher/{email}", handleResearcherTasksList)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
