package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/temporalio/background-checks/types"
)

const ssnTraceAPITimeout = time.Second * 5

func (*Activities) SSNTrace(ctx context.Context, input *types.SSNTraceInput) (*types.SSNTraceResult, error) {
	var result types.SSNTraceResult

	r, err := PostJSON(ctx, "http://thirdparty:8082/ssntrace", input, PostJSONOptions{Timeout: ssnTraceAPITimeout})
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
