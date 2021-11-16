package types

type BackgroundCheckWorkflowInput struct {
	Email   string
	Package string
}

type BackgroundCheckStatus struct {
	Email                      string
	Tier                       string
	Consent                    Consent
	Validate                   ValidateSSNWorkflowResult
	FederalCriminalSearch      FederalCriminalSearchWorkflowResult
	StateCriminalSearch        StateCriminalSearchWorkflowResult
	MotorVehicleIncidentSearch MotorVehicleIncidentSearchWorkflowResult
}

type BackgroundCheckStatusSignal struct {
	ConsentRequired bool
	Status          string
}

type Consent struct {
	Consent  bool
	FullName string
	Address  string
	SSN      string
	DOB      string
}

type ConsentWorkflowResult struct {
	Consent
}

type ConsentRequestSignal struct{}

type ConsentSubmissionSignal struct {
	Consent Consent
}

type ConsentResponseSignal struct {
	Consent Consent
}

type CandidateWorkflowInput struct {
	Email string
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

type ConsentWorkflowInput struct {
	Email string
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
