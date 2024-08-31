package spacex_competitor_flights_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/raphoester/space-trouble-api/internal/pkg/date"
)

func New() *Provider {
	return &Provider{
		httpClient: http.DefaultClient, // might want to inject this later
	}
}

type Provider struct {
	httpClient *http.Client
}

// SpaceX api specific format : UTC, but with milliseconds
const timeFormat = "2006-01-02T15:04:05.000Z"

func (p *Provider) FlightExistsAtLaunchpadOnDate(ctx context.Context, launchpadID string, date date.Date) (bool, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"launchpad": launchpadID,
			"date_utc":  date.Format(timeFormat),
		},
		"options": map[string]interface{}{
			"limit":  1,
			"select": "id",
		},
	}

	marshaled, _ := json.Marshal(query)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.spacexdata.com/v4/launches/query",
		strings.NewReader(string(marshaled)))

	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response dto
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return len(response.Docs) > 0, nil
}

type dto struct {
	Docs []interface{} `json:"docs"`
}
