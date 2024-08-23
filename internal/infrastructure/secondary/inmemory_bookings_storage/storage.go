package inmemory_bookings_storage

import (
	"context"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/pkg/id"
)

func New() *Storage {
	return &Storage{
		bookings: make(map[id.ID]bookings.Booking),
	}
}

type Storage struct {
	bookings map[id.ID]bookings.Booking
}

func (s *Storage) SaveBooking(_ context.Context, booking *bookings.Booking) error {
	s.bookings[booking.ID()] = *booking
	return nil
}

func (s *Storage) ListConflictingFlightBookings(ctx context.Context, booking *bookings.Booking) ([]bookings.Booking, error) {
	ret := make([]bookings.Booking, 0)
	for _, v := range s.bookings {
		if booking.ConflictsWith(v) {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (s *Storage) ListBookings(_ context.Context) ([]bookings.Booking, error) {
	b := make([]bookings.Booking, 0, len(s.bookings))
	for _, booking := range s.bookings {
		b = append(b, booking)
	}
	return b, nil
}
