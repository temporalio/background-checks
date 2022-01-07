package activities

import (
	"io"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func (a *Activities) sendMail(from string, to string, subject string, body io.Reader) error {
	content, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = a.SMTPHost
	server.Port = a.SMTPPort
	server.ConnectTimeout = time.Second
	server.SendTimeout = time.Second

	client, err := server.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject).
		SetBodyData(mail.TextPlain, content)

	if email.Error != nil {
		return email.Error
	}

	return email.Send(client)
}
