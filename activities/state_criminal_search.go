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

const stateCriminalSearchAPITimeout = time.Second * 5

func (*Activities) StateCriminalSearch(ctx context.Context, input *types.StateCriminalSearchInput) (*types.StateCriminalSearchResult, error) {
	var result types.StateCriminalSearchResult

	r, err := PostJSON(ctx, "http://thirdparty:8082/statecriminalsearch", input, PostJSONOptions{Timeout: stateCriminalSearchAPITimeout})
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
