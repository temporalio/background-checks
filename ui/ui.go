package ui

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/utils"
)

const (
	APIEndpoint = "lp-api:8081"
)

type handlers struct{}

//go:embed accept.go.html
var acceptHTML string
var acceptHTMLTemplate = template.Must(template.New("accept").Parse(acceptHTML))

func (h *handlers) handleAccept(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	err := acceptHTMLTemplate.Execute(w, map[string]string{"Token": token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//go:embed accepted.go.html
var acceptedHTML string
var acceptedHTMLTemplate = template.Must(template.New("accepted").Parse(acceptedHTML))

//go:embed declined.go.html
var declinedHTML string
var declinedHTMLTemplate = template.Must(template.New("declined").Parse(declinedHTML))

func (h *handlers) handleAcceptSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	router := api.Router(nil)

	if r.FormValue("action") == "decline" {
		requestURL, err := router.Get("decline").Host(APIEndpoint).URL("token", token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := utils.PostJSON(requestURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		if response.StatusCode != http.StatusOK {
			message := fmt.Sprintf("%s: %s", http.StatusText(response.StatusCode), body)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		err = declinedHTMLTemplate.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	requestURL, err := router.Get("accept").Host(APIEndpoint).URL("token", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	candidatedetails := types.CandidateDetails{
		FullName: r.FormValue("full_name"),
		SSN:      r.FormValue("ssn"),
		Employer: r.FormValue("employer"),
	}
	submission := types.AcceptSubmissionSignal{
		CandidateDetails: candidatedetails,
	}

	response, err := utils.PostJSON(requestURL, submission)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		message := fmt.Sprintf("%s: %s", http.StatusText(response.StatusCode), body)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	err = acceptedHTMLTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Router() *mux.Router {
	r := mux.NewRouter()

	h := handlers{}

	r.HandleFunc("/candidate/{token}", h.handleAccept).Methods("GET")
	r.HandleFunc("/candidate/{token}", h.handleAcceptSubmission).Methods("POST")

	return r
}
