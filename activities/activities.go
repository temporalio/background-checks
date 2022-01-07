package activities

const (
	HiringManagerEmail     = "Hiring Manager <hiring@company.local>"
	HiringSupportEmail     = "BackgroundChecks <support@background-checks.local>"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"
)

type Activities struct {
	SMTPHost string
	SMTPPort int
}
