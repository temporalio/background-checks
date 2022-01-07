package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/types"
)

const declineEmailText = `
{{- $email := .State.Email -}}
{{- $candidate := .State.CandidateDetails -}}
Your background check for: {{$candidate.FullName}} <{{$email}}> has been declined by the candidate.

Thanks,

Background Check System
`

var declineEmailTemplate = template.Must(template.New("declineEmail").Parse(declineEmailText))

func (a *Activities) SendDeclineEmail(ctx context.Context, input types.SendReportEmailInput) (types.SendReportEmailResult, error) {
	var result types.SendReportEmailResult

	var body bytes.Buffer

	err := declineEmailTemplate.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = a.sendMail(HiringSupportEmail, HiringManagerEmail, "Background Check Declined", &body)
	return result, err
}
