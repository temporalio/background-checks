package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/config"
	"github.com/temporalio/background-checks/types"
)

const acceptEmail = `
Hi!

Your potential employer has requested that we run a background check on their behalf.

Please give permission for us proceed with the check by running this command:

"bgc-candidate accept --id {{.CheckID}}"

If you would rather we did not run the check you can decline by running this command:

"bgc-candidate decline --id {{.CheckID}}"

Thanks,

Background Check System
`

func (a *Activities) SendAcceptEmail(ctx context.Context, input types.SendAcceptEmailInput) (types.SendAcceptEmailResult, error) {
	var result types.SendAcceptEmailResult

	var body bytes.Buffer

	t := template.Must(template.New("acceptEmail").Parse(acceptEmail))
	err := t.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = a.SendMail(config.CandidateSupportEmail, input.Email, "Background Check Request", &body)
	if err != nil {
		return result, err
	}

	return result, nil
}
