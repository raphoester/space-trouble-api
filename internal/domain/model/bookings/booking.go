package bookings

import (
	"github.com/raphoester/space-trouble-api/internal/pkg/date"
	"github.com/raphoester/space-trouble-api/internal/pkg/id"
)

func New(id id.ID, clientData ClientData, destinationID string, launchpadID string, launchDate date.Date) *Booking {
	return &Booking{
		id:            id,
		clientData:    clientData,
		destinationID: destinationID,
		launchpadID:   launchpadID,
		launchDate:    launchDate,
	}
}

type Booking struct {
	id            id.ID
	clientData    ClientData
	destinationID string
	launchpadID   string
	launchDate    date.Date
}

type BookingSnapshot struct {
	ID            string
	FirstName     string
	LastName      string
	Gender        string
	Birthday      string
	DestinationID string
	LaunchpadID   string
	LaunchDate    string
}

func (b Booking) ID() id.ID {
	return b.id
}

func (b Booking) ConflictsWith(with Booking) bool {
	if b.launchpadID != with.launchpadID || b.launchDate != with.launchDate {
		return false
	}

	if b.destinationID != with.destinationID {
		return true
	}

	return false
}

func (b Booking) ToSnapshot() BookingSnapshot {
	return BookingSnapshot{
		ID:            b.id.String(),
		FirstName:     b.clientData.FirstName,
		LastName:      b.clientData.LastName,
		Gender:        b.clientData.Gender,
		Birthday:      b.clientData.Birthday.String(),
		DestinationID: b.destinationID,
		LaunchpadID:   b.launchpadID,
		LaunchDate:    b.launchDate.String(),
	}
}
