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

package types

type BackgroundCheckWorkflowInput struct {
	Email   string
	Package string
}

type BackgroundCheckState struct {
	Email            string
	Tier             string
	Accepted         bool
	CandidateDetails CandidateDetails
	SSNTrace         *SSNTraceWorkflowResult
	CheckResults     map[string]interface{}
	CheckErrors      map[string]string
}

type BackgroundCheckWorkflowResult = BackgroundCheckState

type CandidateDetails struct {
	FullName string
	Address  string
	SSN      string
	DOB      string
	Employer string
}

type SendAcceptEmailInput struct {
	Email string
	Token string
}

type SendAcceptEmailResult struct{}

type SendReportEmailInput struct {
	Email string
	State BackgroundCheckState
}

type SendReportEmailResult struct{}

type SendDeclineEmailInput struct {
	Email string
}

type SendDeclineEmailResult struct{}

type AcceptWorkflowInput struct {
	Email string
}

type AcceptWorkflowResult struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

type AcceptSubmission struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

type AcceptSubmissionSignal struct {
	Accepted         bool
	CandidateDetails CandidateDetails
}

type SendEmploymentVerificationEmailInput struct {
	Email string
	Token string
}

type SendEmploymentVerificationEmailResult struct {
}

type EmploymentVerificationWorkflowInput struct {
	CandidateDetails CandidateDetails
}

type EmploymentVerificationWorkflowResult struct {
	EmploymentVerificationComplete bool
	EmployerVerified               bool
}

type EmploymentVerificationSubmission struct {
	EmploymentVerificationComplete bool
	EmployerVerified               bool
}

type EmploymentVerificationSubmissionSignal struct {
	EmploymentVerificationComplete bool
	EmployerVerified               bool
}

type SSNTraceInput struct {
	FullName string
	SSN      string
}
type SSNTraceWorkflowInput struct {
	FullName string
	SSN      string
}

type KnownAddress struct {
	Address string
	City    string
	State   string
	ZipCode string
}

type SSNTraceWorkflowResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

type SSNTraceResult struct {
	SSNIsValid     bool
	KnownAddresses []string
}

type FederalCriminalSearchWorkflowInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchWorkflowResult struct {
	Crimes []string
}

type FederalCriminalSearchInput struct {
	FullName string
	Address  string
}

type FederalCriminalSearchResult struct {
	Crimes []string
}

type StateCriminalSearchWorkflowInput struct {
	FullName       string
	SSNTraceResult []string
}

type StateCriminalSearchWorkflowResult struct {
	Crimes []string
}
type StateCriminalSearchInput struct {
	FullName string
	Address  string
}

type StateCriminalSearchResult struct {
	FullName string
	Address  string
	Crimes   []string
}

type MotorVehicleIncidentSearchWorkflowInput struct {
	FullName string
	Address  string
}

type MotorVehicleIncidentSearchWorkflowResult struct {
	CurrentLicenseState   string
	LicenseValid          bool
	MotorVehicleIncidents []string
}
