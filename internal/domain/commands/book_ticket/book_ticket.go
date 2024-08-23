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

func NewTicketBooker(bookingsRepository BookingsRepository) *TicketBooker {
	return &TicketBooker{
		bookingsRepository: bookingsRepository,
	}
}

type TicketBooker struct {
	bookingsRepository BookingsRepository
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

var ErrLaunchpadUnavailable = errors.New("launchpad is already used for another destination on that day")

func (b *TicketBooker) Execute(
	ctx context.Context,
	params BookTicketParams,
) error {
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

	if err := b.bookingsRepository.SaveBooking(ctx, booking); err != nil {
		return fmt.Errorf("could not save booking: %w", err)
	}

	return nil
}
