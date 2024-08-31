package book_ticket

import (
	"context"
	"errors"
	"fmt"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/pkg/birthday"
	"github.com/raphoester/space-trouble-api/internal/pkg/date"
	"github.com/raphoester/space-trouble-api/internal/pkg/id"
)

func NewTicketBooker(
	bookingsRepository BookingsRepository,
	competitorFlightsProvider CompetitorFlightsProvider,
	launchpadRegistry LaunchpadRegistry,
	destinationRegistry DestinationRegistry,
) *TicketBooker {
	return &TicketBooker{
		bookingsRepository:         bookingsRepository,
		competitorBookingsProvider: competitorFlightsProvider,
		launchpadRegistry:          launchpadRegistry,
		destinationRegistry:        destinationRegistry,
	}
}

type TicketBooker struct {
	bookingsRepository         BookingsRepository
	competitorBookingsProvider CompetitorFlightsProvider
	launchpadRegistry          LaunchpadRegistry
	destinationRegistry        DestinationRegistry
}

type ITicketBooker interface {
	Execute(ctx context.Context, params BookTicketParams) error
}

type BookTicketParams struct {
	ID            string
	FirstName     string
	LastName      string
	Gender        string
	Birthday      string
	LaunchpadID   string
	DestinationID string
	LaunchDate    string
}

type BookingsRepository interface {
	SaveBooking(ctx context.Context, flight *bookings.Booking) error
	ListConflictingFlightBookings(ctx context.Context, booking *bookings.Booking) ([]bookings.Booking, error)
}

type CompetitorFlightsProvider interface {
	FlightExistsAtLaunchpadOnDate(ctx context.Context, launchpadID string, date date.Date) (bool, error)
}

type LaunchpadRegistry interface {
	LaunchpadExists(launchpadID string) bool
}

type DestinationRegistry interface {
	DestinationExists(destinationID string) bool
}

var (
	ErrLaunchpadDoesNotExist   = errors.New("launchpad does not exist")
	ErrDestinationDoesNotExist = errors.New("destination does not exist")
	ErrLaunchpadUnavailable    = errors.New("launchpad is already used for another destination on that day")
)

func (b *TicketBooker) Execute(
	ctx context.Context,
	params BookTicketParams,
) error {
	if !b.launchpadRegistry.LaunchpadExists(params.LaunchpadID) {
		return ErrLaunchpadDoesNotExist
	}

	if !b.destinationRegistry.DestinationExists(params.DestinationID) {
		return ErrDestinationDoesNotExist
	}

	bd, err := birthday.Parse(params.Birthday)
	if err != nil {
		return fmt.Errorf("invalid birthday: %w", err)
	}

	launchDate, err := date.Parse(params.LaunchDate)
	if err != nil {
		return fmt.Errorf("invalid launch date: %w", err)
	}

	booking := bookings.New(id.Parse(params.ID), bookings.ClientData{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Gender:    params.Gender,
		Birthday:  bd,
	}, params.DestinationID, params.LaunchpadID, launchDate)

	conflicts, err := b.bookingsRepository.ListConflictingFlightBookings(ctx, booking)
	if err != nil {
		return fmt.Errorf("could not check launchpad availability: %w", err)
	}

	if len(conflicts) != 0 {
		return ErrLaunchpadUnavailable
	}

	launchpadIsTaken, err := b.competitorBookingsProvider.
		FlightExistsAtLaunchpadOnDate(ctx, params.LaunchpadID, launchDate)
	if err != nil {
		return fmt.Errorf("could not check launchpad availability against provider: %w", err)
	}

	if launchpadIsTaken {
		return ErrLaunchpadUnavailable
	}

	if err := b.bookingsRepository.SaveBooking(ctx, booking); err != nil {
		return fmt.Errorf("could not save booking: %w", err)
	}

	return nil
}
