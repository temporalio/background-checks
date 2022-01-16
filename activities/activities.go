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

package activities

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/temporalio/background-checks/types"
	mail "github.com/xhit/go-simple-mail/v2"
)

const (
	HiringManagerEmail     = "Hiring Manager <hiring@company.local>"
	HiringSupportEmail     = "BackgroundChecks <support@background-checks.local>"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"

	federalCriminalSearchAPITimeout = time.Second * 5
	stateCriminalSearchAPITimeout   = time.Second * 5
	ssnTraceAPITimeout              = time.Second * 5
)

type Activities struct {
	SMTPHost string
	SMTPPort int
	SMTPStub bool
	HTTPStub bool
}

type PostJSONOptions struct {
	Timeout time.Duration
}

func (a *Activities) sendMail(from string, to string, subject string, body io.Reader) error {
	content, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject).
		SetBodyData(mail.TextPlain, content)

	if email.Error != nil {
		return email.Error
	}

	if a.SMTPStub {
		return nil
	}

	server := mail.NewSMTPClient()
	server.Host = a.SMTPHost
	server.Port = a.SMTPPort
	server.ConnectTimeout = time.Second
	server.SendTimeout = time.Second

	client, err := server.Connect()
	if err != nil {
		return err
	}

	return email.Send(client)
}

func (a *Activities) postJSON(ctx context.Context, url string, input interface{}, options PostJSONOptions) (*http.Response, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: options.Timeout,
	}

	return client.Do(req)
}

func (a *Activities) FederalCriminalSearch(ctx context.Context, input *types.FederalCriminalSearchInput) (*types.FederalCriminalSearchResult, error) {
	var result types.FederalCriminalSearchResult

	if a.HTTPStub {
		return &result, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/federalcriminalsearch", input, PostJSONOptions{Timeout: federalCriminalSearchAPITimeout})
	if err != nil {
		return &result, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
}

//go:embed accept_email.go.tmpl
var acceptEmailText string
var acceptEmailTemplate = template.Must(template.New("acceptEmail").Parse(acceptEmailText))

func (a *Activities) SendAcceptEmail(ctx context.Context, input *types.SendAcceptEmailInput) (*types.SendAcceptEmailResult, error) {
	var result types.SendAcceptEmailResult

	var body bytes.Buffer

	err := acceptEmailTemplate.Execute(&body, input)
	if err != nil {
		return &result, err
	}

	err = a.sendMail(CandidateSupportEmail, input.Email, "Background Check Request", &body)
	return &result, err
}

//go:embed decline_email.go.tmpl
var declineEmailText string
var declineEmailTemplate = template.Must(template.New("declineEmail").Parse(declineEmailText))

func (a *Activities) SendDeclineEmail(ctx context.Context, input *types.SendReportEmailInput) (*types.SendReportEmailResult, error) {
	var result types.SendReportEmailResult

	var body bytes.Buffer

	err := declineEmailTemplate.Execute(&body, input)
	if err != nil {
		return &result, err
	}

	err = a.sendMail(HiringSupportEmail, HiringManagerEmail, "Background Check Declined", &body)
	return &result, err
}

//go:embed employment_verification_request.go.tmpl
var employmentVerificationRequestEmailText string
var employmentVerificationRequestEmailTemplate = template.Must(template.New("employmentVerificationRequestEmail").Parse(employmentVerificationRequestEmailText))

func (a *Activities) SendEmploymentVerificationRequestEmail(ctx context.Context, input *types.SendEmploymentVerificationEmailInput) (*types.SendEmploymentVerificationEmailResult, error) {
	var result types.SendEmploymentVerificationEmailResult

	var body bytes.Buffer

	err := employmentVerificationRequestEmailTemplate.Execute(&body, input)
	if err != nil {
		return &result, err
	}

	err = a.sendMail(ResearcherSupportEmail, input.Email, "Employment Verification Request", &body)
	if err != nil {
		return &result, err
	}

	return &result, nil
}

//go:embed report_email.go.tmpl
var reportEmailText string
var reportEmailTemplate = template.Must(template.New("reportEmail").Parse(reportEmailText))

func (a *Activities) SendReportEmail(ctx context.Context, input *types.SendReportEmailInput) (*types.SendReportEmailResult, error) {
	var result types.SendReportEmailResult

	var body bytes.Buffer

	err := reportEmailTemplate.Execute(&body, input.State)
	if err != nil {
		return &result, err
	}

	err = a.sendMail(CandidateSupportEmail, HiringManagerEmail, "Background Check Report", &body)
	return &result, err
}

func (a *Activities) SSNTrace(ctx context.Context, input *types.SSNTraceInput) (*types.SSNTraceResult, error) {
	var result types.SSNTraceResult

	if a.HTTPStub {
		return &types.SSNTraceResult{
			SSNIsValid: true,
		}, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/ssntrace", input, PostJSONOptions{Timeout: ssnTraceAPITimeout})
	if err != nil {
		return &result, err
	}

	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
}

func (a *Activities) StateCriminalSearch(ctx context.Context, input *types.StateCriminalSearchInput) (*types.StateCriminalSearchResult, error) {
	var result types.StateCriminalSearchResult

	if a.HTTPStub {
		return &result, nil
	}

	r, err := a.postJSON(ctx, "http://thirdparty:8082/statecriminalsearch", input, PostJSONOptions{Timeout: stateCriminalSearchAPITimeout})
	if err != nil {
		return &result, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)

		return &result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return &result, err
}
