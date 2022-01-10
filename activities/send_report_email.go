package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/types"
)

const reportEmailText = `
{{- $email := .State.Email -}}
{{- $candidate := .State.CandidateDetails -}}
{{- $ssntrace := .State.SSNTrace -}}
{{- $checks := .State.Checks -}}
Your background check for: {{$candidate.FullName}} <{{$email}}> has completed.

The results are as follows:

SSN Trace:

{{- if $ssntrace.SSNIsValid -}}
SSN is valid
{{- else -}}
**SSN is not valid**
All other checks have been skipped.
{{- end -}}

{{- with $checks.EmploymentVerification -}}
Employment Verification:
{{if .EmployerVerified}}
Verified Employer: {{$candidate.Employer}}
{{else}}
**Employer could not be verified**
{{end}}
{{- end -}}

{{- with $checks.FederalCriminalSearch -}}
Federal Criminal Search:
{{range .Crimes}}
- {{.}}
{{else}}
None found.
{{end}}
{{- end -}}

{{- with $checks.StateCriminalSearch -}}
State Criminal Search:
{{range .Crimes}}
- {{.}}
{{else}}
None found.
{{end}}
{{- end -}}

{{- with $checks.MotorVehicleIncidentSearch -}}
Motor Vehicle Search:

Valid License: {{if .LicenseValid}}Yes{{else}}No{{end}}

Incidents:
{{range .MotorVehicleIncidents}}
- {{.}}
{{else}}
None found.
{{end}}
{{- end -}}

Thanks,

Background Check System
`

var reportEmailTemplate = template.Must(template.New("reportEmail").Parse(reportEmailText))

func (a *Activities) SendReportEmail(ctx context.Context, input *types.SendReportEmailInput) (*types.SendReportEmailResult, error) {
	var result types.SendReportEmailResult

	var body bytes.Buffer

	err := reportEmailTemplate.Execute(&body, input)
	if err != nil {
		return &result, err
	}

	err = a.sendMail(CandidateSupportEmail, HiringManagerEmail, "Background Check Report", &body)
	return &result, err
}
