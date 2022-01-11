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

const federalCriminalSearchAPITimeout = time.Second * 5

func (a *Activities) FederalCriminalSearch(ctx context.Context, input *types.FederalCriminalSearchInput) (*types.FederalCriminalSearchResult, error) {
	var result types.FederalCriminalSearchResult

	r, err := a.postJSON(ctx, "http://thirdparty:8082/federalcriminalsearch", input, PostJSONOptions{Timeout: federalCriminalSearchAPITimeout})
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
