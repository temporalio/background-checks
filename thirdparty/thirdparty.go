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
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temporalio/background-checks/types"
)

const DefaultEndpoint = "localhost:8082"

func handleSsnTrace(w http.ResponseWriter, r *http.Request) {
	var input types.ValidateSSNInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"SSN" : "123-45-6789", "Address" : "123 Main St", "FullName" : "Joe Bloggs"}

	var result types.ValidateSSNResult
	result.Valid = false
	if input.SSN == "123-45-6789" {
		result.Valid = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func handleMotorVehicleSearch(w http.ResponseWriter, r *http.Request) {
	var input types.MotorVehicleIncidentSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.MotorVehicleIncidentSearchResult
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

func handleEmploymentSearch(w http.ResponseWriter, r *http.Request) {
	var input types.EmploymentSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.EmploymentSearchResult
	if input.FullName == "Joe Bloggs" {
		companies := []string{"Buy n Large, Inc.", "Acme, Co."}
		result.Companies = companies
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func handleFederalCriminalSearch(w http.ResponseWriter, r *http.Request) {
	var input types.FederalCriminalSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.FederalCriminalSearchResult
	if input.FullName == "Joe Bloggs" {
		crimes := []string{"Money Laundering", "Pick-pocketing"}
		result.Crimes = crimes
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func handleStateCriminalSearch(w http.ResponseWriter, r *http.Request) {
	var input types.StateCriminalSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test Input - {"FullName" : "Joe Bloggs", "Address" : "123 Main St"}

	var result types.StateCriminalSearchResult
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
	r.HandleFunc("/employmentsearch/", handleEmploymentSearch).Name("employmentsearch")
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
