package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func PostJSON(url *url.URL, input interface{}) (*http.Response, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	return client.Do(req)
}

func GetJSON(url *url.URL, result interface{}) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %v", err)
	}

	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 200 && r.StatusCode < 300 {
		err = json.NewDecoder(r.Body).Decode(result)
		return r, err
	}

	defer r.Body.Close()
	message, _ := ioutil.ReadAll(r.Body)

	return r, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), message)
}
