package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"

	"github.com/github/go-fault"
	"github.com/gorilla/mux"
)

const DefaultEndpoint = "0.0.0.0:8082"

// For convenience this fake third party API happens to take the same shape inputs and return the same
// shape results as our application uses.

type SSNTraceInput struct {
	FullName string
	SSN      string
}

type SSNTraceResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

type MotorVehicleIncidentSearchInput struct {
	FullName string
	Address  string
}

type MotorVehicleIncidentSearchResult struct {
	LicenseValid          bool
	MotorVehicleIncidents []string
}

type FederalCriminalSearchInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchResult struct {
	Crimes []string
}

type StateCriminalSearchInput struct {
	FullName string
	Address  string
}

type StateCriminalSearchResult struct {
	FullName string
	Address  string
	Crimes   []string
}

func handleSsnTrace(w http.ResponseWriter, r *http.Request) {
	var input SSNTraceInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result SSNTraceResult

	var validSSN = regexp.MustCompile(`\d{3}-\d{2}-\d{4}$`)
	result.SSNIsValid = validSSN.MatchString(input.SSN)

	addressMap := map[string][]string{
		"111-11-1111": {"123 Broadway, New York, NY 10011", "1 E. 161 St, Bronx, NY 10451", "41 Seaver Way, Queens, NY 11368"},
		"222-22-2222": {"456 Oak Street, Springfield, IL 62706", "1060 W. Addison St, Chicago, IL 60613"},
		"333-33-3333": {"4 Jersey St, Boston, MA 02215", "333 W Camden St, Baltimore, MD 21201"},
		"444-44-4444": {"1 Royal Way, Kansas City, MO 64129", "", "700 Clark Ave, St Louis, MO 63102"}}

	result.KnownAddresses = addressMap[input.SSN]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleMotorVehicleSearch(w http.ResponseWriter, r *http.Request) {
	var input MotorVehicleIncidentSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result MotorVehicleIncidentSearchResult
	var motorVehicleIncidents []string

	possibleMotorVehicleIncidents := []string{
		"Speeding",
		"Reckless Driving",
		"Driving Without Insurance",
		"Driving Under the Influence",
	}

	rndnum := rand.Intn(100)
	if rndnum > 25 {
		motorVehicleIncidents = append(motorVehicleIncidents, possibleMotorVehicleIncidents[rand.Intn(len(possibleMotorVehicleIncidents))])
		result.LicenseValid = false
	} else {
		result.LicenseValid = true
	}

	result.MotorVehicleIncidents = motorVehicleIncidents

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleFederalCriminalSearch(w http.ResponseWriter, r *http.Request) {
	var input FederalCriminalSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result FederalCriminalSearchResult
	var crimes []string

	possibleCrimes := []string{
		"Money Laundering",
		"Racketeering",
		"Counterfeiting",
		"Espionage",
	}

	rndnum := rand.Intn(100)
	if rndnum > 75 {
		crimes = append(crimes, possibleCrimes[rand.Intn(len(possibleCrimes))])
	}
	result.Crimes = crimes

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleStateCriminalSearch(w http.ResponseWriter, r *http.Request) {
	var input StateCriminalSearchInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result StateCriminalSearchResult
	var crimes []string

	possibleCrimes := []string{
		"Jaywalking",
		"Burglary",
		"Foul Play",
		"Smoking in a Public Place",
	}

	rndnum := rand.Intn(100)
	if rndnum > 75 {
		crimes = append(crimes, possibleCrimes[rand.Intn(len(possibleCrimes))])
	}
	result.FullName = input.FullName
	result.Address = input.Address
	result.Crimes = crimes

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
