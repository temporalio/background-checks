package config

import "time"

const (
	TaskQueue                 = "background-checks-main"
	SMTPServer                = "lp-mailhog:1025"
	CandidateSupportEmail     = "candidates@background-checks.local"
	ResearcherSupportEmail    = "BackgroundChecks <researchers@background-checks.local>"
	CandidateSupportEmailOld  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmailOld = "BackgroundChecks <researchers@background-checks.local>"
	AcceptGracePeriod         = time.Hour * 24 * 7
)
