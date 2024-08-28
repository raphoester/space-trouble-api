package controller

import (
	"context"
	"fmt"

	bookingsv1 "github.com/raphoester/space-trouble-api/generated/proto/bookings/v1"
	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
)

func New(ticketBooker *book_ticket.TicketBooker) *Controller {
	return &Controller{
		ticketBooker: ticketBooker,
	}
}

type Controller struct {
	bookingsv1.UnimplementedBookingsServiceServer
	ticketBooker *book_ticket.TicketBooker
}

func (c *Controller) BookTicket(ctx context.Context,
	req *bookingsv1.BookTicketRequest) (*bookingsv1.BookTicketResponse, error) {

	err := c.ticketBooker.Execute(ctx, book_ticket.BookTicketParams{
		ID:            req.GetId(),
		FirstName:     req.GetFirstName(),
		LastName:      req.GetLastName(),
		Gender:        req.GetGender(),
		Birthday:      req.GetBirthday(),
		LaunchpadID:   req.GetLaunchpadId(),
		DestinationID: req.GetDestinationId(),
		LaunchDate:    req.GetLaunchDate(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to book ticket: %w", err)
	}

	return &bookingsv1.BookTicketResponse{}, nil
}
