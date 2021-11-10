package types

type BackgroundCheckInput struct {
	Email string
	Tier  string
}

type ConsentInput struct {
	Email string
}

type ConsentResult struct {
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

type SearchResult struct {
	Type                             string
	FederalCriminalSearchResult      FederalCriminalSearchResult
	StateCriminalSearchResult        StateCriminalSearchResult
	MotorVehicleIncidentSearchResult MotorVehicleIncidentSearchResult
}

func (r SearchResult) Result() interface{} {
	switch r.Type {
	case "FederalCriminalSearchResult":
		return r.FederalCriminalSearchResult
	case "StateCriminalSearchResult":
		return r.StateCriminalSearchResult
	case "MotorVehicleIncidentSearchResult":
		return r.MotorVehicleIncidentSearchResult
	}

	return nil
}
