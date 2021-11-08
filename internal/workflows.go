package internal

type AcceptCheckInput struct {
	Email string
}

type AcceptCheckResult struct {
	Accept   bool
	FullName string
	Address  string
	SSN      string
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
