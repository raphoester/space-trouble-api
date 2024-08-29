package controller

import (
	"context"
	"fmt"

	bookingsv1 "github.com/raphoester/space-trouble-api/generated/proto/bookings/v1"
	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings"
)

func New(
	ticketBooker book_ticket.ITicketBooker,
	bookingsGetter get_all_bookings.Getter,
) *Controller {
	return &Controller{
		ticketBooker:   ticketBooker,
		bookingsGetter: bookingsGetter,
	}
}

type Controller struct {
	bookingsv1.UnimplementedBookingsServiceServer
	ticketBooker   book_ticket.ITicketBooker
	bookingsGetter get_all_bookings.Getter
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

func (c *Controller) GetAllBookings(ctx context.Context,
	_ *bookingsv1.GetAllBookingsRequest) (*bookingsv1.GetAllBookingsResponse, error) {

	bookings, err := c.bookingsGetter.GetAllBookings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting bookings: %w", err)
	}

	protoBookings := make([]*bookingsv1.Booking, 0, len(bookings))
	for _, b := range bookings {
		protoBookings = append(protoBookings, &bookingsv1.Booking{
			Id:            b.Id,
			FirstName:     b.FirstName,
			LastName:      b.LastName,
			LaunchpadId:   b.LaunchpadID,
			DestinationId: b.DestinationID,
			LaunchDate:    b.LaunchDate,
		})
	}

	return &bookingsv1.GetAllBookingsResponse{
		Bookings: protoBookings,
	}, nil
}
