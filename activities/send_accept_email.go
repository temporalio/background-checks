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

const acceptEmail = `To: {{.Email}}
Subject: Background Check Request

Hi!

Your potential employer has requested that we run a background check on their behalf.

Please give permission for us proceed with the check by running this command:

"bgc-candidate accept --id {{.CheckID}}"

If you would rather we did not run the check you can decline by running this command:

"bgc-candidate decline --id {{.CheckID}}"

Thanks,

Background Check System
`

func SendAcceptEmail(ctx context.Context, input types.SendAcceptEmailInput) (types.SendAcceptEmailResult, error) {
	var result types.SendAcceptEmailResult

	to := []string{input.Email}
	var body bytes.Buffer

	fmt.Fprintf(&body, "From: %s\n", config.CandidateSupportEmail)

	t := template.Must(template.New("acceptEmail").Parse(acceptEmail))
	err := t.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = smtp.SendMail(config.SMTPServer, nil, config.CandidateSupportEmail, to, body.Bytes())
	if err != nil {
		return result, err
	}

	return result, nil
}
