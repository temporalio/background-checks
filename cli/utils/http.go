package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func PostJSON(url *url.URL, input interface{}) (*http.Response, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request input: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	return client.Do(req)
}
