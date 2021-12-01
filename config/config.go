package config

import "time"

const (
	TaskQueue              = "background-checks-main"
	SMTPServer             = "lp-mailhog:1025"
	HiringManagerEmail     = "Hiring Manager <hiring@company.local"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"
	AcceptGracePeriod      = time.Hour * 24 * 7
	ResearchDeadline       = time.Hour * 24 * 7
	ThirdPartyAPIEndpoint  = "lp-thirdparty-api:8082"
)
