package psql_bookings_storage

import (
	"context"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
)

type Storage struct {
	postgres *postgres.Postgres
}

func (s *Storage) SaveBooking(ctx context.Context, booking *bookings.Booking) error {
	if err := s.postgres.Gorm.WithContext(ctx).Save(booking.ToSnapshot()).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListConflictingFlightBookings(_ context.Context, booking *bookings.Booking) ([]bookings.Booking, error) {
	panic("not implemented")
}
