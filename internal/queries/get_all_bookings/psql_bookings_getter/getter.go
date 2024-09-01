package psql_bookings_getter

import (
	"context"
	"fmt"

	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings"
)

func New(pg *postgres.Postgres) *Getter {
	return &Getter{
		pg: pg,
	}
}

type Getter struct {
	pg *postgres.Postgres
}

func (g *Getter) GetAllBookings(ctx context.Context) ([]get_all_bookings.Booking, error) {
	rows, err := g.pg.DB.QueryContext(ctx,
		"SELECT id, first_name, last_name, launchpad_id, destination_id, launch_date FROM bookings")
	if err != nil {
		return nil, fmt.Errorf("failed to query db: %w", err)
	}

	bookings := make([]get_all_bookings.Booking, 0, 5)
	for rows.Next() {
		var b get_all_bookings.Booking
		if err := rows.Scan(&b.Id, &b.FirstName, &b.LastName, &b.LaunchpadID, &b.DestinationID, &b.LaunchDate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}
