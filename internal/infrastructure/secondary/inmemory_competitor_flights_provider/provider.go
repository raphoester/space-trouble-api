package inmemory_competitor_flights_provider

import (
	"context"
	"errors"
	"sync"

	"github.com/raphoester/space-trouble-api/internal/pkg/date"
)

func New() *Provider {
	return &Provider{
		flights: make([]flight, 0),
	}
}

type Provider struct {
	mu      sync.RWMutex
	flights []flight
}

type flight struct {
	launchpadID string
	date        date.Date
}

func (p *Provider) FlightExistsAtLaunchpadOnDate(ctx context.Context, launchpadID string,
	departureDate date.Date) (bool, error) {

	return p.flightExists(launchpadID, departureDate), nil
}

func (p *Provider) RegisterFlight(launchpadID string, departureDate date.Date) error {
	if p.flightExists(launchpadID, departureDate) {
		return errors.New("launchpad is unavailable at that date")
	}

	p.mu.Lock()
	p.flights = append(p.flights, flight{
		launchpadID: launchpadID,
		date:        departureDate,
	})
	p.mu.Unlock()
	return nil
}

func (p *Provider) flightExists(launchpadID string, departureDate date.Date) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, f := range p.flights {
		if f.date == departureDate && f.launchpadID == launchpadID {
			return true
		}
	}
	return false
}
