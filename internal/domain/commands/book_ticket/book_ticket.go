package book_ticket

import (
	"context"
	"fmt"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/pkg/birthday"
	"github.com/raphoester/space-trouble-api/internal/pkg/date"
	"github.com/raphoester/space-trouble-api/internal/pkg/id"
)

func NewBookTicket() *BookTicket {
	return &BookTicket{}
}

type BookTicket struct {
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
}

func (b *BookTicket) Execute(
	ctx context.Context,
	params BookTicketParams,
	bookingsRepository BookingsRepository,
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

	if err := bookingsRepository.SaveBooking(ctx, booking); err != nil {
		return fmt.Errorf("could not save booking: %w", err)
	}

	return nil
}
