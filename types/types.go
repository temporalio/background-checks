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
	Validate                   ValidateSSNWorkflowResult
	FederalCriminalSearch      FederalCriminalSearchWorkflowResult
	StateCriminalSearch        StateCriminalSearchWorkflowResult
	MotorVehicleIncidentSearch MotorVehicleIncidentSearchWorkflowResult
}

type CandidateDetails struct {
	FullName string
	Address  string
	SSN      string
	DOB      string
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

type AcceptWorkflowInput struct {
}

type ValidateSSNWorkflowInput struct {
	FullName string
	Address  string
	SSN      string
}

type ValidateSSNWorkflowResult struct {
	Valid bool
}

type FederalCriminalSearchWorkflowInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchWorkflowResult struct {
	Crimes []string
}

type StateCriminalSearchWorkflowInput struct {
	FullName string
	Address  string
}

type StateCriminalSearchWorkflowResult struct {
	Crimes []string
}

type MotorVehicleIncidentSearchWorkflowInput struct {
	FullName string
	Address  string
}

type EmploymentSearchWorkflowInput struct {
	FullName string
	Address  string
}

type EmploymentSearchWorkflowResult struct {
	Companies []string
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
