package activities

import (
	"bytes"
	"context"
	"text/template"

	"github.com/temporalio/background-checks/types"
)

const acceptEmailText = `
Hello, 

Your potential employer has requested that we conduct a background check on their behalf.

The following information is needed to complete your check:
- Full Name (as it appears on your government ID)
- Social Security Number
- Current Employer

Please give permission for us proceed with the check by running this command and adding your details:

./run-cli bgc-candidate accept --token {{.Token}} --fullname 'your name here' --ssn '111-11-1111' --employer 'Your current employer'

If you would rather we did not run the check you can decline by running this command:

./run-cli bgc-candidate decline --token {{.Token}}

Thanks,

Background Check System
`

var acceptEmailTemplate = template.Must(template.New("acceptEmail").Parse(acceptEmailText))

func (a *Activities) SendAcceptEmail(ctx context.Context, input types.SendAcceptEmailInput) (types.SendAcceptEmailResult, error) {
	var result types.SendAcceptEmailResult

	var body bytes.Buffer

	err := acceptEmailTemplate.Execute(&body, input)
	if err != nil {
		return result, err
	}

	err = a.sendMail(CandidateSupportEmail, input.Email, "Background Check Request", &body)
	return result, err
}
