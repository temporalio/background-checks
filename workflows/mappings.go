package workflows

import (
	"encoding/base64"
	"fmt"
	"path"

	"go.temporal.io/sdk/workflow"
)

func BackgroundCheckWorkflowID(email string) string {
	return fmt.Sprintf("BackgroundCheck:%s", email)
}

func AcceptWorkflowID(email string) string {
	return fmt.Sprintf("Accept:%s", email)
}

func EmploymentVerificationWorkflowID(email string) string {
	return fmt.Sprintf("EmploymentVerification:%s", email)
}

func CheckWorkflowID(email string, name string) string {
	return fmt.Sprintf("%s:%s", name, email)
}

func TokenForWorkflow(ctx workflow.Context) string {
	info := workflow.GetInfo(ctx)

	rawToken := path.Join(info.WorkflowExecution.ID, info.WorkflowExecution.RunID)

	return base64.URLEncoding.EncodeToString([]byte(rawToken))
}

func WorkflowFromToken(token string) (string, string, error) {
	var rawToken []byte

	rawToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", "", err
	}

	wfid := path.Dir(string(rawToken))
	runid := path.Base(string(rawToken))

	return wfid, runid, nil
}
