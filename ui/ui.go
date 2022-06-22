package ui

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/temporalio/background-checks/api"
	"github.com/temporalio/background-checks/utils"
	"github.com/temporalio/background-checks/workflows"
)

const (
	APIEndpoint = "api:8081"
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

	candidatedetails := workflows.CandidateDetails{
		FullName: r.FormValue("full_name"),
		SSN:      r.FormValue("ssn"),
		Employer: r.FormValue("employer"),
	}
	submission := workflows.AcceptSubmissionSignal{
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

//go:embed employment_verification.go.html
var employmentVerificationHTML string
var employmentVerificationHTMLTemplate = template.Must(template.New("employment_verification").Parse(employmentVerificationHTML))

func (h *handlers) handleEmploymentVerification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	router := api.Router(nil)

	requestURL, err := router.Get("employmentverify_details").Host(APIEndpoint).URL("token", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var candidate workflows.CandidateDetails

	_, err = utils.GetJSON(requestURL, &candidate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = employmentVerificationHTMLTemplate.Execute(w, map[string]interface{}{"Token": token, "Candidate": candidate})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//go:embed employment_verified.go.html
var employmentVerifiedHTML string
var employmentVerifiedHTMLTemplate = template.Must(template.New("employment_verification").Parse(employmentVerifiedHTML))

func (h *handlers) handleEmploymentVerificationSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	router := api.Router(nil)

	requestURL, err := router.Get("employmentverify").Host(APIEndpoint).URL("token", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	submission := workflows.EmploymentVerificationSubmissionSignal{
		EmploymentVerificationComplete: true,
		EmployerVerified:               r.FormValue("action") == "yes",
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

	err = employmentVerifiedHTMLTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//go:embed status.go.html
var statusHTML string
var statusHTMLTemplate = template.Must(template.New("status").Parse(statusHTML))

func (h *handlers) handleStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	router := api.Router(nil)

	requestURL, err := router.Get("check_status").Host(APIEndpoint).URL("token", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var status workflows.BackgroundCheckState

	_, err = utils.GetJSON(requestURL, &status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = statusHTMLTemplate.Execute(w, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//go:embed report.go.html
var reportHTML string
var reportHTMLTemplate = template.Must(template.New("report").Parse(reportHTML))

func (h *handlers) handleReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	router := api.Router(nil)

	requestURL, err := router.Get("check_status").Host(APIEndpoint).URL("token", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var status workflows.BackgroundCheckState

	_, err = utils.GetJSON(requestURL, &status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = reportHTMLTemplate.Execute(w, status)
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

	r.HandleFunc("/employment/{token}", h.handleEmploymentVerification).Methods("GET")
	r.HandleFunc("/employment/{token}", h.handleEmploymentVerificationSubmission).Methods("POST")

	r.HandleFunc("/status/{token}", h.handleStatus).Methods("GET")
	r.HandleFunc("/report/{token}", h.handleReport).Methods("GET")

	return r
}
