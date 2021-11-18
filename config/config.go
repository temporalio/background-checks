package config

import "time"

const (
	TaskQueue              = "background-checks-main"
	SMTPServer             = "localhost:1025"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"
	AcceptGracePeriod      = time.Hour * 24 * 7
)
