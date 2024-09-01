package bookings

import (
	"fmt"

	"github.com/raphoester/space-trouble-api/internal/pkg/birthday"
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

func Restore(snapshot *BookingSnapshot) (*Booking, error) {
	anID := id.Parse(snapshot.ID)

	launchDate, err := date.Parse(snapshot.LaunchDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse launch date: %w", err)
	}

	bd, err := birthday.Parse(snapshot.Birthday)
	if err != nil {
		return nil, fmt.Errorf("failed to parse birthday: %w", err)
	}

	return &Booking{
		id: anID,
		clientData: ClientData{
			FirstName: snapshot.FirstName,
			LastName:  snapshot.LastName,
			Gender:    snapshot.Gender,
			Birthday:  bd,
		},
		destinationID: snapshot.DestinationID,
		launchpadID:   snapshot.LaunchpadID,
		launchDate:    launchDate,
	}, nil
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
