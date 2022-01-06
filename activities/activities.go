package activities

import "net/smtp"

type Activities struct {
	SMTPServer string
	SMTPAuth   smtp.Auth
}
