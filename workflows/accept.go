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
	"time"

	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/background-checks/types"
)

const (
	AcceptSubmissionSignal = "accept-submission"
	AcceptGracePeriod      = time.Hour * 24 * 7
)

func emailCandidate(ctx workflow.Context, input *types.AcceptWorkflowInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	})

	i := types.SendAcceptEmailInput{
		Email: input.Email,
		Token: TokenForWorkflow(ctx),
	}
	f := workflow.ExecuteActivity(ctx, a.SendAcceptEmail, i)
	return f.Get(ctx, nil)
}

func waitForSubmission(ctx workflow.Context) (*types.AcceptSubmission, error) {
	var response types.AcceptSubmission
	var err error

	s := workflow.NewSelector(ctx)

	ch := workflow.GetSignalChannel(ctx, AcceptSubmissionSignal)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
		var submission types.AcceptSubmissionSignal
		c.Receive(ctx, &submission)

		response = types.AcceptSubmission(submission)
	})
	s.AddFuture(workflow.NewTimer(ctx, AcceptGracePeriod), func(f workflow.Future) {
		err = f.Get(ctx, nil)

		// Treat failure to accept in time as declining.
		response.Accepted = false
	})

	s.Select(ctx)

	return &response, err
}

// @@@SNIPSTART background-checks-accept-workflow-definition
func Accept(ctx workflow.Context, input *types.AcceptWorkflowInput) (*types.AcceptWorkflowResult, error) {
	err := emailCandidate(ctx, input)
	if err != nil {
		return &types.AcceptWorkflowResult{}, err
	}

	submission, err := waitForSubmission(ctx)

	result := types.AcceptWorkflowResult(*submission)
	return &result, err
}

// @@@SNIPEND
