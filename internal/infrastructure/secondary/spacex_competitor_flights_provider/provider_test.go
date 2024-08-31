package spacex_competitor_flights_provider_test

import (
	"context"
	"testing"

	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/spacex_competitor_flights_provider"
	"github.com/raphoester/space-trouble-api/internal/pkg/date"
	"github.com/stretchr/testify/require"
)

func TestProvider_FlightExists(t *testing.T) {
	provider := spacex_competitor_flights_provider.New()
	res, err := provider.FlightExistsAtLaunchpadOnDate(context.Background(),
		"5e9e4501f509094ba4566f84", date.MustParse("01/12/2022"))

	require.NoError(t, err)
	require.True(t, res)
}

func TestProvider_FlightDoesNotExists(t *testing.T) {
	provider := spacex_competitor_flights_provider.New()
	res, err := provider.FlightExistsAtLaunchpadOnDate(context.Background(),
		"5e9e4501f509094ba4566f84", date.MustParse("14/07/2022")) // no flight on this date
	require.NoError(t, err)
	require.False(t, res)
}
