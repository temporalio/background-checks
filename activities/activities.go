package activities

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

const (
	HiringManagerEmail     = "Hiring Manager <hiring@company.local>"
	HiringSupportEmail     = "BackgroundChecks <support@background-checks.local>"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"

	federalCriminalSearchAPITimeout = time.Second * 5
	stateCriminalSearchAPITimeout   = time.Second * 5
	ssnTraceAPITimeout              = time.Second * 5
)

type Activities struct {
	SMTPHost string
	SMTPPort int
	SMTPStub bool
	HTTPStub bool
}

type PostJSONOptions struct {
	Timeout time.Duration
}

func (a *Activities) sendMail(from string, to string, subject string, htmlTemplate *template.Template, textTemplate *template.Template, input interface{}) error {
	var htmlContent, textContent bytes.Buffer

	err := htmlTemplate.Execute(&htmlContent, input)
	if err != nil {
		return err
	}

	err = textTemplate.Execute(&textContent, input)
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject).
		SetBody(mail.TextHTML, htmlContent.String()).
		AddAlternative(mail.TextPlain, textContent.String())

	if email.Error != nil {
		return email.Error
	}

	if a.SMTPStub {
		return nil
	}

	server := mail.NewSMTPClient()
	server.Host = a.SMTPHost
	server.Port = a.SMTPPort
	server.ConnectTimeout = time.Second
	server.SendTimeout = time.Second

	client, err := server.Connect()
	if err != nil {
		return err
	}

	return email.Send(client)
}

func (a *Activities) postJSON(ctx context.Context, url string, input interface{}, options PostJSONOptions) (*http.Response, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: options.Timeout,
	}

	return client.Do(req)
}

type FederalCriminalSearchInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchResult struct {
	Crimes []string
}

func (a *Activities) FederalCriminalSearch(ctx context.Context, input *FederalCriminalSearchInput) (*FederalCriminalSearchResult, error) {
	var result FederalCriminalSearchResult

	if a.HTTPStub {
		return &result, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/federalcriminalsearch", input, PostJSONOptions{Timeout: federalCriminalSearchAPITimeout})
	if err != nil {
		return &result, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
}

//go:embed accept_email.go.html
var acceptEmailHTML string
var acceptEmailHTMLTemplate = template.Must(template.New("acceptEmailHTML").Parse(acceptEmailHTML))

//go:embed accept_email.go.tmpl
var acceptEmailText string
var acceptEmailTextTemplate = template.Must(template.New("acceptEmailText").Parse(acceptEmailText))

type SendAcceptEmailInput struct {
	Email string
	Token string
}

type SendAcceptEmailResult struct{}

func (a *Activities) SendAcceptEmail(ctx context.Context, input *SendAcceptEmailInput) (*SendAcceptEmailResult, error) {
	var result SendAcceptEmailResult

	err := a.sendMail(CandidateSupportEmail, input.Email, "Background Check Request", acceptEmailHTMLTemplate, acceptEmailTextTemplate, input)
	return &result, err
}

//go:embed decline_email.go.html
var declineEmailHTML string
var declineEmailHTMLTemplate = template.Must(template.New("declineEmailHTML").Parse(declineEmailHTML))

//go:embed decline_email.go.tmpl
var declineEmailText string
var declineEmailTextTemplate = template.Must(template.New("declineEmailText").Parse(declineEmailText))

type SendDeclineEmailInput struct {
	Email string
}

type SendDeclineEmailResult struct{}

func (a *Activities) SendDeclineEmail(ctx context.Context, input *SendDeclineEmailInput) (*SendDeclineEmailResult, error) {
	var result SendDeclineEmailResult

	err := a.sendMail(HiringSupportEmail, HiringManagerEmail, "Background Check Declined", declineEmailHTMLTemplate, declineEmailTextTemplate, input)
	return &result, err
}

//go:embed employment_verification_request.go.html
var employmentVerificationRequestEmailHTML string
var employmentVerificationRequestEmailHTMLTemplate = template.Must(template.New("employmentVerificationRequestEmailHTML").Parse(employmentVerificationRequestEmailHTML))

//go:embed employment_verification_request.go.tmpl
var employmentVerificationRequestEmailText string
var employmentVerificationRequestEmailTextTemplate = template.Must(template.New("employmentVerificationRequestEmailText").Parse(employmentVerificationRequestEmailText))

type SendEmploymentVerificationEmailInput struct {
	Email string
	Token string
}

type SendEmploymentVerificationEmailResult struct{}

func (a *Activities) SendEmploymentVerificationRequestEmail(ctx context.Context, input *SendEmploymentVerificationEmailInput) (*SendEmploymentVerificationEmailResult, error) {
	var result SendEmploymentVerificationEmailResult

	err := a.sendMail(ResearcherSupportEmail, input.Email, "Employment Verification Request", employmentVerificationRequestEmailHTMLTemplate, employmentVerificationRequestEmailTextTemplate, input)

	return &result, err
}

//go:embed report_email.go.html
var reportEmailHTML string
var reportEmailHTMLTemplate = template.Must(template.New("reportEmailHTML").Parse(reportEmailHTML))

//go:embed report_email.go.tmpl
var reportEmailText string
var reportEmailTextTemplate = template.Must(template.New("reportEmailText").Parse(reportEmailText))

type SendReportEmailInput struct {
	Email string
	Token string
}

type SendReportEmailResult struct{}

func (a *Activities) SendReportEmail(ctx context.Context, input *SendReportEmailInput) (*SendReportEmailResult, error) {
	var result SendReportEmailResult

	err := a.sendMail(CandidateSupportEmail, HiringManagerEmail, "Background Check Report", reportEmailHTMLTemplate, reportEmailTextTemplate, input)
	return &result, err
}

type SSNTraceInput struct {
	FullName string
	SSN      string
}

type SSNTraceResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

func (a *Activities) SSNTrace(ctx context.Context, input *SSNTraceInput) (*SSNTraceResult, error) {
	var result SSNTraceResult

	if a.HTTPStub {
		return &SSNTraceResult{
			SSNIsValid: true,
		}, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/ssntrace", input, PostJSONOptions{Timeout: ssnTraceAPITimeout})
	if err != nil {
		return &result, err
	}

	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
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

func (a *Activities) StateCriminalSearch(ctx context.Context, input *StateCriminalSearchInput) (*StateCriminalSearchResult, error) {
	var result StateCriminalSearchResult

	if a.HTTPStub {
		return &result, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/statecriminalsearch", input, PostJSONOptions{Timeout: stateCriminalSearchAPITimeout})
	if err != nil {
		return &result, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
}

type MotorVehicleIncidentSearchInput struct {
	FullName string
	Address  string
}

type MotorVehicleIncidentSearchResult struct {
	LicenseValid          bool
	MotorVehicleIncidents []string
}

func (a *Activities) MotorVehicleIncidentSearch(ctx context.Context, input *MotorVehicleIncidentSearchInput) (*MotorVehicleIncidentSearchResult, error) {
	var result MotorVehicleIncidentSearchResult

	if a.HTTPStub {
		return &result, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/motorvehiclesearch", input, PostJSONOptions{Timeout: stateCriminalSearchAPITimeout})
	if err != nil {
		return &result, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
}
