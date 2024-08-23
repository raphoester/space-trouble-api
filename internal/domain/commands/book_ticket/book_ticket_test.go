package book_ticket_test

import (
	"context"
	"testing"

	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/pkg/date"
	"github.com/raphoester/space-trouble-api/internal/pkg/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite
}

func (s *testSuite) TestNominalCase() {
	s.T().Run("Should be able to book a ticket when no conflict is detected", func(t *testing.T) {
		storage := inmemory_bookings_storage.New()
		bookTicket := book_ticket.NewTicketBooker(storage)

		bookingID := (&id.FixedIDGenerator{}).Generate().String()

		params := book_ticket.BookTicketParams{
			ID:            bookingID,
			FirstName:     "John",
			LastName:      "Doe",
			Gender:        "Male",
			Birthday:      "13/05",
			LaunchpadID:   "example-launchpad-id",
			DestinationID: "example-destination-id",
			LaunchDate:    "13/10/2024",
		}

		err := bookTicket.Execute(context.Background(), params)

		assert.NoError(t, err)
		savedBookings, err := storage.ListBookings(context.Background())
		assert.NoError(t, err)
		require.Len(t, savedBookings, 1)

		expectedSnapshot := bookings.BookingSnapshot{
			ID:            bookingID,
			FirstName:     "John",
			LastName:      "Doe",
			Gender:        "Male",
			Birthday:      "13/05",
			DestinationID: "example-destination-id",
			LaunchpadID:   "example-launchpad-id",
			LaunchDate:    "13/10/2024",
		}

		assert.EqualValues(t, expectedSnapshot, savedBookings[0].ToSnapshot())
	})
}

func (s *testSuite) TestConflict() {
	s.T().Run("Should not be able to book a ticket when the booking triggers a conflict", func(t *testing.T) {
		storage := inmemory_bookings_storage.New()
		bookTicket := book_ticket.NewTicketBooker(storage)

		idFactory := id.NewChaoticFactory(t.Name())
		launchpadID := idFactory.Generate()
		destinationID := idFactory.Generate()
		launchDate := date.MustParse("25/11/2024")

		conflictingBooking := bookings.New((&id.FixedIDGenerator{}).Generate(),
			bookings.ClientData{
				FirstName: "John",
				LastName:  "Doe",
				Gender:    "Male",
				Birthday:  "07/03",
			}, destinationID.String(), launchpadID.String(), launchDate)
		err := storage.SaveBooking(context.Background(), conflictingBooking)
		require.NoError(t, err)

		err = bookTicket.Execute(context.Background(),
			book_ticket.BookTicketParams{
				ID:            idFactory.Generate().String(),
				FirstName:     "Jane",
				LastName:      "Smith",
				Gender:        "Female",
				Birthday:      "29/09",
				LaunchpadID:   launchpadID.String(),
				DestinationID: idFactory.Generate().String(), // other destination
				LaunchDate:    launchDate.String(),
			})

		assert.Error(t, err)
		assert.ErrorIs(t, err, book_ticket.ErrLaunchpadUnavailable)

		storedBookings, err := storage.ListBookings(context.Background())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(storedBookings), 1)
		assert.Len(t, storedBookings, 1)

		firstBooking := storedBookings[0]
		assert.EqualValues(t, conflictingBooking.ToSnapshot(), firstBooking.ToSnapshot())
	})
}
