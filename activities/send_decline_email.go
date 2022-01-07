package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/types"
)

const declineEmail = `
{{- $email := .State.Email -}}
{{- $candidate := .State.CandidateDetails -}}
Your background check for: {{$candidate.FullName}} <{{$email}}> has been declined by the candidate.

Thanks,

Background Check System
`

func (a *Activities) SendDeclineEmail(ctx context.Context, input types.SendReportEmailInput) (types.SendReportEmailResult, error) {
	var result types.SendReportEmailResult

	var body bytes.Buffer

	t := template.Must(template.New("declineEmail").Parse(declineEmail))
	err := t.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = a.sendMail(config.HiringSupportEmail, config.HiringManagerEmail, "Background Check Declined", &body)
	return result, err
}
