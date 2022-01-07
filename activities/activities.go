package activities

import "net/smtp"

const (
	SMTPServer             = "lp-mailhog:1025"
	HiringManagerEmail     = "Hiring Manager <hiring@company.local"
	HiringSupportEmail     = "BackgroundChecks <support@background-checks.local>"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"
)

type Activities struct {
	SMTPServer string
	SMTPAuth   smtp.Auth
}
