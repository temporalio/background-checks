package activities

import (
	"bytes"
	"fmt"
	"io"
	"net/smtp"

	"github.com/temporalio/background-checks/config"
)

type Activities struct {
	Server string
	Auth   smtp.Auth
}

func (a *Activities) SendMail(from string, to string, subject string, body io.Reader) error {
	var b bytes.Buffer

	fmt.Fprintf(&b, "From: %s\nTo: %s\nSubject: %s\n\n", from, to, subject)

	_, err := io.Copy(&b, body)
	if err != nil {
		return err
	}

	return smtp.SendMail(a.Server, a.Auth, config.CandidateSupportEmail, []string{to}, b.Bytes())
}
