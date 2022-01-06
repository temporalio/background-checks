package mappings

import (
	"fmt"
)

func BackgroundCheckWorkflowID(email string) string {
	return fmt.Sprintf("BackgroundCheck:%s", email)
}

func AcceptWorkflowID(checkID string) string {
	return fmt.Sprintf("Accept:%s", checkID)
}

func EmploymentVerificationWorkflowID(checkID string) string {
	return fmt.Sprintf("EmploymentVerification:%s", checkID)
}

func CheckWorkflowID(checkID string, name string) string {
	return fmt.Sprintf("%s:%s", name, checkID)
}
