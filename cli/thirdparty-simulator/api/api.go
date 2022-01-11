package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/github/go-fault"
	"github.com/gorilla/mux"

	"github.com/temporalio/background-checks/types"
)

const DefaultEndpoint = "0.0.0.0:8082"

func handleSsnTrace(w http.ResponseWriter, r *http.Request) {
	var input types.SSNTraceWorkflowInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result types.SSNTraceWorkflowResult
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
	var err error

	errorInjector, _ := fault.NewErrorInjector(http.StatusInternalServerError)
	errorFault, _ := fault.NewFault(errorInjector,
		fault.WithEnabled(true),
		fault.WithParticipation(0.3),
	)

	handlerChain := errorFault.Handler(Router())

	srv := &http.Server{
		Handler: handlerChain,
		Addr:    DefaultEndpoint,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		srv.Close()
	case err = <-errCh:
		log.Fatalf("error: %v", err)
	}
}
