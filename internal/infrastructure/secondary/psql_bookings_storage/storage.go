package psql_bookings_storage

import (
	"context"
	"fmt"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/pkg/id"
	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
)

func New(pg *postgres.Postgres) *Storage {
	return &Storage{
		postgres: pg,
	}
}

type Storage struct {
	postgres *postgres.Postgres
}

func (s *Storage) SaveBooking(ctx context.Context, booking *bookings.Booking) error {
	if err := s.postgres.Gorm.WithContext(ctx).Table("bookings").Save(booking.ToSnapshot()).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListConflictingFlightBookings(ctx context.Context, booking *bookings.Booking) ([]*bookings.Booking, error) {
	snapshot := booking.ToSnapshot()
	matches := make([]*bookings.BookingSnapshot, 0, 3)
	driver := s.postgres.Gorm.WithContext(ctx).
		Table("bookings").
		Where(`launchpad_id = ? AND launch_date = ? AND destination_id != ?`,
			snapshot.LaunchpadID, snapshot.LaunchDate, snapshot.DestinationID).
		Find(&matches)

	if driver.Error != nil {
		return nil, fmt.Errorf("failed to list conflicting flight bookings: %w", driver.Error)
	}

	restored := make([]*bookings.Booking, 0, len(matches))
	for _, match := range matches {
		restoredBooking, err := bookings.Restore(match)
		if err != nil {
			return nil, fmt.Errorf("failed to restore snapshot: %w", err)
		}
		restored = append(restored, restoredBooking)
	}

	return restored, nil
}

func (s *Storage) Get(ctx context.Context, id id.ID) (*bookings.Booking, error) {
	var snapshot bookings.BookingSnapshot
	if err := s.postgres.Gorm.Table("bookings").WithContext(ctx).First(&snapshot, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	restored, err := bookings.Restore(&snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to restore snapshot: %w", err)
	}

	return restored, nil
}
