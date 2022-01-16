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
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/types"
)

const (
	EmploymentVerificationDetailsQuery     = "employment-verification-details"
	EmploymentVerificationSubmissionSignal = "employment-verification-submission"
	ResearchDeadline                       = time.Hour * 24 * 7
)

func chooseResearcher(ctx workflow.Context, input *types.EmploymentVerificationWorkflowInput) (string, error) {
	researchers := []string{
		"researcher1@example.com",
		"researcher2@example.com",
		"researcher3@example.com",
	}

	// Here we just pick a random researcher.
	// In a real use case you may round-robin, decide based on price or current workload,
	// or fetch a researcher from a third party API.

	var researcher string
	r := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return researchers[rand.Intn(len(researchers))]
	})
	err := r.Get(&researcher)

	return researcher, err
}

func emailEmploymentVerificationRequest(ctx workflow.Context, input *types.EmploymentVerificationWorkflowInput, email string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	evsend := workflow.ExecuteActivity(ctx, a.SendEmploymentVerificationRequestEmail, types.SendEmploymentVerificationEmailInput{
		Email: email,
		Token: TokenForWorkflow(ctx),
	})
	return evsend.Get(ctx, nil)
}

func waitForEmploymentVerificationSubmission(ctx workflow.Context) (*types.EmploymentVerificationSubmission, error) {
	var response types.EmploymentVerificationSubmission
	var err error

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, EmploymentVerificationSubmissionSignal)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.EmploymentVerificationSubmissionSignal
		c.Receive(ctx, &submission)

		response = types.EmploymentVerificationSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, ResearchDeadline), func(f workflow.Future) {
		err = f.Get(ctx, nil)

		// We should probably fail the (child) workflow here.
		response.EmploymentVerificationComplete = false
		response.EmployerVerified = false
	})

	s.Select(ctx)

	return &response, err
}

// @@@SNIPSTART background-checks-employment-verification-workflow-definition
func EmploymentVerification(ctx workflow.Context, input *types.EmploymentVerificationWorkflowInput) (*types.EmploymentVerificationWorkflowResult, error) {
	var result types.EmploymentVerificationWorkflowResult

	err := workflow.SetQueryHandler(ctx, EmploymentVerificationDetailsQuery, func() (types.CandidateDetails, error) {
		return input.CandidateDetails, nil
	})
	if err != nil {
		return &result, err
	}

	researcher, err := chooseResearcher(ctx, input)
	if err != nil {
		return &result, err
	}

	err = emailEmploymentVerificationRequest(ctx, input, researcher)
	if err != nil {
		return &result, err
	}
	submission, err := waitForEmploymentVerificationSubmission(ctx)

	result = types.EmploymentVerificationWorkflowResult(*submission)
	return &result, err
}

// @@@SNIPEND
