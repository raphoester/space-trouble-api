package get_all_bookings

import "context"

type Booking struct {
	Id            string
	FirstName     string
	LastName      string
	LaunchpadID   string
	DestinationID string
	LaunchDate    string
}

type Getter interface {
	GetAllBookings(ctx context.Context) ([]Booking, error)
}
