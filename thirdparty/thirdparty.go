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
package thirdparty

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sdkclient "go.temporal.io/sdk/client"

	"github.com/gorilla/mux"

	"github.com/temporalio/background-checks/mappings"
	"github.com/temporalio/background-checks/signals"
	"github.com/temporalio/background-checks/types"
)

const DefaultEndpoint = "localhost:8082"

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

func handleSsnTrace(w http.ResponseWriter, r *http.Request) {
	var input types.ValidateSSNWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"SSN" : "123-45-6789", "Address" : "123 Main St", "FullName" : "Joe Bloggs"}

	var result types.ValidateSSNWorkflowResult
	result.Valid = false
	if input.SSN == "123-45-6789" {
		result.Valid = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func handleMotorVehicleSearch(w http.ResponseWriter, r *http.Request) {
	var input types.MotorVehicleIncidentSearchWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.MotorVehicleIncidentSearchWorkflowResult
	result.LicenseValid = false
	if input.FullName == "Joe Bloggs" {
		result.LicenseValid = true
		result.CurrentLicenseState = "CA"
		incident := []string{"License Revoked 12/1/2020"}
		result.MotorVehicleIncidents = incident
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func handleEmploymentVerification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input types.CandidateDetails
	log.Println("Employment Verification Entry: ", input.Employer)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St", "DOB" : "1/1/1990", "Employer" : "Acme Corp"}

	var result types.EmploymentVerificationWorkflowResult

	result.CandidateDetails = input

	err = signalWorkflow(
		mappings.EmploymentVerificationWorkflowID(id),
		signals.EmploymentVerificationSubmission,
		result,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/* if verified {
		result.CandidateDetails.EmployerVerified = verified
		return
	} */
	result.EmployerVerificationComplete = true

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func signalWorkflow(wid string, signalName string, signalArg interface{}) error {
	c, err := getClient()
	if err != nil {
		return err
	}

	return c.SignalWorkflow(context.Background(), wid, "", signalName, signalArg)
}

func handleFederalCriminalSearch(w http.ResponseWriter, r *http.Request) {
	var input types.FederalCriminalSearchWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.FederalCriminalSearchWorkflowResult
	if input.FullName == "Joe Bloggs" {
		crimes := []string{"Money Laundering", "Pick-pocketing"}
		result.Crimes = crimes
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func handleStateCriminalSearch(w http.ResponseWriter, r *http.Request) {
	var input types.StateCriminalSearchWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.StateCriminalSearchWorkflowResult
	if input.FullName == "Joe Bloggs" {
		crimes := []string{"Jay-walking", "Littering"}
		result.Crimes = crimes
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ssntrace/", handleSsnTrace).Methods("POST").Name("ssntrace")
	r.HandleFunc("/motorvehiclesearch/", handleMotorVehicleSearch).Methods("POST").Name("motorvehiclesearch")
	r.HandleFunc("/employmentverification/", handleEmploymentVerification).Name("employmentverification")
	r.HandleFunc("/federalcriminalsearch/", handleFederalCriminalSearch).Methods("POST").Name("federalcriminalsearch")
	r.HandleFunc("/statecriminalsearch/", handleStateCriminalSearch).Name("statecriminalsearch")
	return r
}

func Run() {
	srv := &http.Server{
		Handler: Router(),
		Addr:    DefaultEndpoint,
	}

	log.Fatal(srv.ListenAndServe())
}
