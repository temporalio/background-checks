package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/temporalio/background-checks/types"
)

const stateCriminalSearchAPITimeout = time.Second * 5

func (a *Activities) StateCriminalSearch(ctx context.Context, input types.StateCriminalSearchInput) (types.StateCriminalSearchResult, error) {
	var result types.StateCriminalSearchResult

	requestURL, err := url.Parse("http://thirdparty:8082/statecriminalsearch")
	if err != nil {
		return result, err
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return result, fmt.Errorf("unable to encode input: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(jsonInput))
	if err != nil {
		return result, fmt.Errorf("unable to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: stateCriminalSearchAPITimeout,
	}
	r, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		return result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return result, err
}
