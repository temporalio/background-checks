package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/types"
)

const employmentVerificationRequestEmail = `
{{ $candidate := .CandidateDetails -}}
Hello Background Check Researcher, 

Our candidate {{$candidate.FullName}} is undergoing a background check, and the next step is to verify their employment history. 

Please reach out to their company and confirm whether they are currently employed. 

Candidate Name: {{$candidate.FullName}}
Employer: {{$candidate.Employer}}

When you have completed this step, respond by updating the Background Check using the instructions below:

EMPLOYMENT IS VERIFIED:

./run-cli bgc-researcher employmentverify --token {{.Token}}

EMPLOYMENT IS NOT VERIFIED:

TBA

Thanks,

Background Check System
`

func (a *Activities) SendEmploymentVerificationRequestEmail(ctx context.Context, input types.SendEmploymentVerificationEmailInput) (types.SendEmploymentVerificationEmailResult, error) {
	var result types.SendEmploymentVerificationEmailResult

	var body bytes.Buffer

	t := template.Must(template.New("employmentVerificationRequestEmail").Parse(employmentVerificationRequestEmail))
	err := t.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = a.sendMail(config.ResearcherSupportEmail, input.Email, "Employment Verification Request", &body)
	if err != nil {
		return result, err
	}

	return result, nil
}
