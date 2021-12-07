package types

type BackgroundCheckWorkflowInput struct {
	Email   string
	Package string
}

type BackgroundCheckStatus int64

const (
	BackgroundCheckStatusUnknown BackgroundCheckStatus = iota
	BackgroundCheckStatusPendingAccept
	BackgroundCheckStatusRunning
	BackgroundCheckStatusCompleted
	BackgroundCheckStatusDeclined
	BackgroundCheckStatusFailed
	BackgroundCheckStatusTerminated
	BackgroundCheckStatusCancelled
)

func (s BackgroundCheckStatus) String() string {
	switch s {
	case BackgroundCheckStatusPendingAccept:
		return "pending_consent"
	case BackgroundCheckStatusRunning:
		return "running"
	case BackgroundCheckStatusCompleted:
		return "completed"
	case BackgroundCheckStatusDeclined:
		return "declined"
	case BackgroundCheckStatusFailed:
		return "failed"
	case BackgroundCheckStatusTerminated:
		return "terminated"
	case BackgroundCheckStatusCancelled:
		return "cancelled"
	}

	return "unknown"
}

func BackgroundCheckStatusFromString(s string) BackgroundCheckStatus {
	switch s {
	case BackgroundCheckStatusPendingAccept.String():
		return BackgroundCheckStatusPendingAccept
	case BackgroundCheckStatusRunning.String():
		return BackgroundCheckStatusRunning
	case BackgroundCheckStatusCompleted.String():
		return BackgroundCheckStatusCompleted
	case BackgroundCheckStatusDeclined.String():
		return BackgroundCheckStatusDeclined
	case BackgroundCheckStatusFailed.String():
		return BackgroundCheckStatusFailed
	case BackgroundCheckStatusTerminated.String():
		return BackgroundCheckStatusTerminated
	case BackgroundCheckStatusCancelled.String():
		return BackgroundCheckStatusCancelled
	default:
		return BackgroundCheckStatusUnknown
	}
}

type BackgroundCheckState struct {
	Email                      string
	Tier                       string
	Accepted                   bool
	CandidateDetails           CandidateDetails
	ValidateSSN                ValidateSSNWorkflowResult
	EmploymentVerification     EmploymentVerificationWorkflowResult
	FederalCriminalSearch      FederalCriminalSearchWorkflowResult
	StateCriminalSearch        StateCriminalSearchWorkflowResult
	MotorVehicleIncidentSearch MotorVehicleIncidentSearchWorkflowResult
}

type CandidateDetails struct {
	FullName         string
	Address          string
	SSN              string
	DOB              string
	Employer         string
	EmployerVerified bool
}

type SendAcceptEmailInput struct {
	Email   string
	CheckID string
}

type SendAcceptEmailResult struct{}

type SendReportEmailInput struct {
	Email string
	State BackgroundCheckState
}

type SendReportEmailResult struct{}

type AcceptWorkflowInput struct {
	Email   string
	CheckID string
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
	Email            string
	CandidateDetails CandidateDetails
	CheckID          string
}

type SendEmploymentVerificationEmailResult struct {
}

type EmploymentVerificationWorkflowInput struct {
	CandidateDetails CandidateDetails
	CheckID          string
}

type EmploymentVerificationWorkflowResult struct {
	CandidateDetails             CandidateDetails
	EmployerVerificationComplete bool
}

type EmploymentVerificationSubmission struct {
	CandidateDetails             CandidateDetails
	EmployerVerificationComplete bool
}

type EmploymentVerificationSubmissionSignal struct {
	CandidateDetails             CandidateDetails
	EmployerVerificationComplete bool
}

type ResearcherWorkflowInput struct {
	Email string
}

type ResearcherTodo struct {
	Token                           string
	Type                            string
	FederalCriminalSearchInput      FederalCriminalSearchWorkflowInput
	StateCriminalSearchInput        StateCriminalSearchWorkflowInput
	MotorVehicleIncidentSearchInput MotorVehicleIncidentSearchWorkflowInput
}

func (r ResearcherTodo) Input() interface{} {
	switch r.Type {
	case "FederalCriminalSearch":
		return r.FederalCriminalSearchInput
	case "StateCriminalSearch":
		return r.StateCriminalSearchInput
	case "MotorVehicleIncidentSearch":
		return r.MotorVehicleIncidentSearchInput
	}

	return nil
}

type SSNTraceInput struct {
	FullName string
	SSN      string
}
type ValidateSSNWorkflowInput struct {
	FullName string
	SSN      string
}

type KnownAddress struct {
	Address string
	City    string
	State   string
	ZipCode string
}

type ValidateSSNWorkflowResult struct {
	KnownAddresses []string
}

type SSNTraceResult struct {
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

type SearchResult struct {
	Type                             string
	FederalCriminalSearchResult      FederalCriminalSearchWorkflowResult
	StateCriminalSearchResult        StateCriminalSearchWorkflowResult
	MotorVehicleIncidentSearchResult MotorVehicleIncidentSearchWorkflowResult
}

func (r SearchResult) Result() interface{} {
	switch r.Type {
	case "FederalCriminalSearch":
		return r.FederalCriminalSearchResult
	case "StateCriminalSearch":
		return r.StateCriminalSearchResult
	case "MotorVehicleIncidentSearch":
		return r.MotorVehicleIncidentSearchResult
	}

	return nil
}
