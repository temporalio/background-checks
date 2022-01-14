package types

type BackgroundCheckWorkflowInput struct {
	Email   string
	Package string
}

type BackgroundCheckState struct {
	Email            string
	Tier             string
	Accepted         bool
	CandidateDetails CandidateDetails
	SSNTrace         *SSNTraceWorkflowResult
	CheckResults     map[string]interface{}
	CheckErrors      map[string]string
}

type BackgroundCheckWorkflowResult = BackgroundCheckState

type CandidateDetails struct {
	FullName string
	Address  string
	SSN      string
	DOB      string
	Employer string
}

type SendAcceptEmailInput struct {
	Email string
	Token string
}

type SendAcceptEmailResult struct{}

type SendReportEmailInput struct {
	Email string
	State BackgroundCheckState
}

type SendReportEmailResult struct{}

type SendDeclineEmailInput struct {
	Email string
}

type SendDeclineEmailResult struct{}

type AcceptWorkflowInput struct {
	Email string
}

type AcceptWorkflowResult struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

type AcceptSubmission struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

type AcceptSubmissionSignal struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

type SendEmploymentVerificationEmailInput struct {
	Email string
	Token string
}

type SendEmploymentVerificationEmailResult struct {
}

type EmploymentVerificationWorkflowInput struct {
	CandidateDetails CandidateDetails
}

type EmploymentVerificationWorkflowResult struct {
	EmploymentVerificationComplete bool
	EmployerVerified               bool
}

type EmploymentVerificationSubmission struct {
	EmploymentVerificationComplete bool
	EmployerVerified               bool
}

type EmploymentVerificationSubmissionSignal struct {
	EmploymentVerificationComplete bool
	EmployerVerified               bool
}

type SSNTraceInput struct {
	FullName string
	SSN      string
}
type SSNTraceWorkflowInput struct {
	FullName string
	SSN      string
}

type KnownAddress struct {
	Address string
	City    string
	State   string
	ZipCode string
}

type SSNTraceWorkflowResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

type SSNTraceResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

type FederalCriminalSearchWorkflowInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchWorkflowResult struct {
	Crimes []string
}

type FederalCriminalSearchInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchResult struct {
	Crimes []string
}

type StateCriminalSearchWorkflowInput struct {
	FullName       string
	SSNTraceResult []string
}

type StateCriminalSearchWorkflowResult struct {
	Crimes []string
}
type StateCriminalSearchInput struct {
	FullName string
	Address  string
}

type StateCriminalSearchResult struct {
	FullName string
	Address  string
	Crimes   []string
}

type MotorVehicleIncidentSearchWorkflowInput struct {
	FullName string
	Address  string
}

type MotorVehicleIncidentSearchWorkflowResult struct {
	CurrentLicenseState   string
	LicenseValid          bool
	MotorVehicleIncidents []string
}
