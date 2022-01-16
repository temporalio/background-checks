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

package mocks

import "github.com/temporalio/background-checks/types"

var SSNTraceWorkflowResult = map[types.SSNTraceWorkflowInput]types.SSNTraceWorkflowResult{
	{FullName: "John Smith", SSN: "111-11-1111"}:    {},
	{FullName: "Sally Jones", SSN: "123-45-6789"}:   {},
	{FullName: "Javier Bardem", SSN: "987-65-4321"}: {},
}

/*
{
	{Address: "123 Broadway", City: "New York", State: "NY", ZipCode: "10011"},
	{Address: "500 Market Street", City: "San Francisco", State: "CA", ZipCode: "94110"},
	{Address: "111 Dearborn Ave", City: "Detroit", State: "MI", ZipCode: "44014"},
} */
