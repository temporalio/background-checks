package mappings

import "fmt"

func BackgroundCheckWorkflowID(email string) string {
	return fmt.Sprintf("BackgroundCheck-%s", email)
}

func CandidateWorkflowID(email string) string {
	return fmt.Sprintf("Candidate-%s", email)
}

func ResearcherWorkflowID(email string) string {
	return fmt.Sprintf("Researcher-%s", email)
}
