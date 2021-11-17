package mappings

import "fmt"

func BackgroundCheckWorkflowID(email string) string {
	return fmt.Sprintf("BackgroundCheck:%s", email)
}

func AcceptWorkflowID(checkID string) string {
	return fmt.Sprintf("Accept:%s", checkID)
}

func ResearcherWorkflowID(email string) string {
	return fmt.Sprintf("Researcher:%s", email)
}
