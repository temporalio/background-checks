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
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/temporalio/background-checks/types"
)

const DefaultEndpoint = "0.0.0.0:8082"

func handleSsnTrace(w http.ResponseWriter, r *http.Request) {
	var input types.ValidateSSNWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result types.ValidateSSNWorkflowResult
	if input.SSN == "111-11-1111" {
		Addresses := []string{
			"123 Broadway, New York, NY 10011",
			"500 Market Street, San Francisco, CA 94110",
			"111 Dearborn Ave, Detroit, MI 44014",
		}
		result.SSNIsValid = true
		result.KnownAddresses = Addresses
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
	if input.FullName == "John Smith" {
		result.LicenseValid = true
		result.CurrentLicenseState = "CA"
		incident := []string{"License Revoked 12/1/2020"}
		result.MotorVehicleIncidents = incident
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

	var result types.FederalCriminalSearchResult
	if input.FullName == "John Smith" {
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

	var result types.StateCriminalSearchResult
	result.FullName = input.FullName
	result.Address = input.Address

	if input.Address == "500 Market Street, San Francisco, CA 94110" && input.FullName == "John Smith" {
		crimes := []string{"Jay-walking", "Littering"}
		result.Crimes = crimes
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ssntrace", handleSsnTrace).Methods("POST")
	r.HandleFunc("/motorvehiclesearch", handleMotorVehicleSearch).Methods("POST")
	r.HandleFunc("/federalcriminalsearch", handleFederalCriminalSearch).Methods("POST")
	r.HandleFunc("/statecriminalsearch", handleStateCriminalSearch).Methods("POST")

	return r
}

func Run() {
	srv := &http.Server{
		Handler: Router(),
		Addr:    DefaultEndpoint,
	}

	log.Fatal(srv.ListenAndServe())
}
