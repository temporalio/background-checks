package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/types"
)

const reportEmail = `
Your background check for: {{.State.CandidateDetails.FullName}} <{{.State.Email}}> has completed.

The results are as follows:

Employment Verification:
{{if .State.EmploymentVerification.EmploymentVerificationComplete}}
Verified: {{if .State.EmploymentVerification.EmployerVerified}}Yes{{else}}No{{end}}
{{else}}
Not completed.
{{end}}

Federal Criminal Search:
{{range .State.FederalCriminalSearch.Crimes}}
- {{.}}
{{else}}
None found.
{{end}}

State Criminal Search:
{{range .State.StateCriminalSearch.Crimes}}
- {{.}}
{{else}}
None found.
{{end}}

Motor Vehicle Search:

Valid License: {{if .State.MotorVehicleIncidentSearch.LicenseValid}}Yes{{else}}No{{end}}

Incidents:
{{range .State.MotorVehicleIncidentSearch.MotorVehicleIncidents}}
- {{.}}
{{else}}
None found.
{{end}}

Employment Verification:
{{if .State.EmploymentVerification.EmployerVerified}}Verified Employer: {{.State.CandidateDetails.Employer}}{{else}}**Employer could not be verified**{{end}}

Thanks,

Background Check System
`

func (a *Activities) SendReportEmail(ctx context.Context, input types.SendReportEmailInput) (types.SendReportEmailResult, error) {
	var result types.SendReportEmailResult

	var body bytes.Buffer

	t := template.Must(template.New("reportEmail").Parse(reportEmail))
	err := t.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = a.SendMail(config.HiringSupportEmail, config.HiringManagerEmail, "Background Check Report", &body)
	if err != nil {
		return result, err
	}

	return result, nil
}
