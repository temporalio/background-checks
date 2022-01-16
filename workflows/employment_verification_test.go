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

package workflows_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/background-checks/activities"
	"github.com/temporalio/background-checks/types"
	"github.com/temporalio/background-checks/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestEmploymentVerificationWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	details := types.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.OnActivity(a.SendEmploymentVerificationRequestEmail, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input *types.SendEmploymentVerificationEmailInput) (*types.SendEmploymentVerificationEmailResult, error) {
			return &types.SendEmploymentVerificationEmailResult{}, nil
		},
	)

	env.RegisterDelayedCallback(
		func() {
			env.SignalWorkflow(
				workflows.EmploymentVerificationSubmissionSignal,
				types.EmploymentVerificationSubmissionSignal{EmploymentVerificationComplete: true, EmployerVerified: true},
			)
		},
		0,
	)

	env.ExecuteWorkflow(workflows.EmploymentVerification, &types.EmploymentVerificationWorkflowInput{CandidateDetails: details})

	var result types.EmploymentVerificationWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, types.EmploymentVerificationWorkflowResult{EmploymentVerificationComplete: true, EmployerVerified: true}, result)
}

func TestEmploymentVerificationWorkflowTimeout(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	details := types.CandidateDetails{
		FullName: "John Smith",
		SSN:      "111-11-1111",
		DOB:      "1981-01-01",
		Address:  "1 Chestnut Avenue",
	}

	env.OnActivity(a.SendEmploymentVerificationRequestEmail, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input *types.SendEmploymentVerificationEmailInput) (*types.SendEmploymentVerificationEmailResult, error) {
			return &types.SendEmploymentVerificationEmailResult{}, nil
		},
	)

	env.ExecuteWorkflow(workflows.EmploymentVerification, &types.EmploymentVerificationWorkflowInput{CandidateDetails: details})

	var result types.EmploymentVerificationWorkflowResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, types.EmploymentVerificationWorkflowResult{EmploymentVerificationComplete: false, EmployerVerified: false}, result)
}
