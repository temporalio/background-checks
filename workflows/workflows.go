/*
 * The MIT License
 *
 * Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
 *
 * Copyright (c) 2020 Uber Technologies, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

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
