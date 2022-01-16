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

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func PostJSON(url *url.URL, input interface{}) (*http.Response, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	return client.Do(req)
}

func GetJSON(url *url.URL, result interface{}) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode >= 200 && r.StatusCode < 300 {
		err = json.NewDecoder(r.Body).Decode(result)
		return r, err
	}

	message, _ := io.ReadAll(r.Body)

	return r, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), message)
}
