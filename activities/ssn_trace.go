package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/temporalio/background-checks/types"
)

const ssnTraceAPITimeout = time.Second * 5

func (*Activities) SSNTrace(ctx context.Context, input types.SSNTraceInput) (types.SSNTraceResult, error) {
	var result types.SSNTraceResult

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return result, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://thirdparty:8082/ssntrace", bytes.NewReader(jsonInput))
	if err != nil {
		return result, fmt.Errorf("unable to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: ssnTraceAPITimeout,
	}
	r, err := client.Do(req)
	if err != nil {
		return result, err
	}

	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		return result, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), body)
	}

	err = json.NewDecoder(r.Body).Decode(&result)
	return result, err
}
