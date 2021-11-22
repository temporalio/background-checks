package activities

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/types"
)

const employmentVerificationRequestEmail = `To: Researcher
Subject: Employment Verification Request

Hello Background Check Researcher, 

Our candidate is undergoing a background check, and the next step is to verify their employment history. 

Please reach out to their company and confirm whether they are currently employed. 

When you have completed this step, respond by updating the Background Check using the instructions below:

EMPLOYMENT IS VERIFIED:
"bgc-researcher employmentverify --id {{.CheckID}}"

EMPLOYMENT IS NOT VERIFIED:
TBA

Thanks,

Background Check System
`

func SendEmploymentVerificationRequestEmail(ctx context.Context, input types.SendEmploymentVerificationEmailInput) (types.SendAcceptEmailResult, error) {
	var result types.SendAcceptEmailResult

	var to = []string{config.ResearcherSupportEmail}
	var body bytes.Buffer

	fmt.Fprintf(&body, "From: %s\n", config.CandidateSupportEmail)

	t := template.Must(template.New("employmentVerificationRequestEmail").Parse(employmentVerificationRequestEmail))
	err := t.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = smtp.SendMail(config.SMTPServer, nil, config.ResearcherSupportEmail, to, body.Bytes())
	if err != nil {
		return result, err
	}

	return result, nil
}
