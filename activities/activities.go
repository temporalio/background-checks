package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	HiringManagerEmail     = "Hiring Manager <hiring@company.local>"
	HiringSupportEmail     = "BackgroundChecks <support@background-checks.local>"
	CandidateSupportEmail  = "BackgroundChecks <candidates@background-checks.local>"
	ResearcherSupportEmail = "BackgroundChecks <researchers@background-checks.local>"
)

type Activities struct {
	SMTPHost string
	SMTPPort int
}

type PostJSONOptions struct {
	Timeout time.Duration
}

func PostJSON(ctx context.Context, url string, input interface{}, options PostJSONOptions) (*http.Response, error) {
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
