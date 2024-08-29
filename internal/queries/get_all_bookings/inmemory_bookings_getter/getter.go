package inmemory_bookings_getter

import (
	"context"

	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings"
)

func New(storage *inmemory_bookings_storage.Storage) *Getter {
	return &Getter{
		storage: storage,
	}
}

type Getter struct {
	storage *inmemory_bookings_storage.Storage
}

func (g *Getter) GetAllBookings(_ context.Context) ([]get_all_bookings.Booking, error) {
	bookings := g.storage.ListBookings()
	ret := make([]get_all_bookings.Booking, 0, len(bookings))
	for _, b := range bookings {
		snapshot := b.ToSnapshot()
		ret = append(ret, get_all_bookings.Booking{
			Id:            snapshot.ID,
			FirstName:     snapshot.FirstName,
			LastName:      snapshot.LastName,
			LaunchpadID:   snapshot.LaunchpadID,
			DestinationID: snapshot.DestinationID,
			LaunchDate:    snapshot.LaunchDate,
		})
	}
	return ret, nil
}
