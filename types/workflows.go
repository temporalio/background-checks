package types

type BackgroundCheckInput struct {
	Email string
	Tier  string
}

type CandidateInput struct {
	Email string
}

type CandidateTodo struct {
	Token string
}

type ResearcherInput struct {
	Email string
}

type ResearcherTodo struct {
	Token                           string
	Type                            string
	FederalCriminalSearchInput      FederalCriminalSearchInput
	StateCriminalSearchInput        StateCriminalSearchInput
	MotorVehicleIncidentSearchInput MotorVehicleIncidentSearchInput
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

type ConsentInput struct {
	Email string
}

type ConsentResult struct {
	Consent  bool
	FullName string
	Address  string
	SSN      string
	DOB      string
}

type ValidateSSNInput struct {
	FullName string
	Address  string
	SSN      string
}

type ValidateSSNResult struct {
	Valid bool
}

type FederalCriminalSearchInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchResult struct {
	Crimes []string
}

type StateCriminalSearchInput struct {
	FullName string
	Address  string
}

type StateCriminalSearchResult struct {
	Crimes []string
}

type MotorVehicleIncidentSearchInput struct {
	FullName string
	Address  string
}

type MotorVehicleIncidentSearchResult struct {
	Crimes []string
}

type SearchResult struct {
	Type                             string
	FederalCriminalSearchResult      FederalCriminalSearchResult
	StateCriminalSearchResult        StateCriminalSearchResult
	MotorVehicleIncidentSearchResult MotorVehicleIncidentSearchResult
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