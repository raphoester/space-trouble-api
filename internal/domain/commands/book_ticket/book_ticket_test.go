package book_ticket_test

import (
	"context"
	"testing"

	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_launchpad_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_competitor_flights_provider"
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
		competitorFlightsProvider := inmemory_competitor_flights_provider.New()
		launchpadRegistry := hardcoded_launchpad_registry.New()
		bookTicket := book_ticket.NewTicketBooker(storage, competitorFlightsProvider, launchpadRegistry)

		bookingID := (&id.FixedIDGenerator{}).Generate().String()

		params := book_ticket.BookTicketParams{
			ID:            bookingID,
			FirstName:     "John",
			LastName:      "Doe",
			Gender:        "Male",
			Birthday:      "13/05",
			LaunchpadID:   hardcoded_launchpad_registry.Florida,
			DestinationID: "example-destination-id",
			LaunchDate:    "13/10/2024",
		}

		err := bookTicket.Execute(context.Background(), params)

		assert.NoError(t, err)
		savedBookings := storage.ListBookings()
		require.Len(t, savedBookings, 1)

		expectedSnapshot := bookings.BookingSnapshot{
			ID:            bookingID,
			FirstName:     "John",
			LastName:      "Doe",
			Gender:        "Male",
			Birthday:      "13/05",
			DestinationID: "example-destination-id",
			LaunchpadID:   hardcoded_launchpad_registry.Florida,
			LaunchDate:    "13/10/2024",
		}

		assert.EqualValues(t, expectedSnapshot, savedBookings[0].ToSnapshot())
	})
}

func (s *testSuite) TestInternalConflict() {
	s.T().Run("Should not be able to book a ticket when the booking conflicts "+
		"with another previous booking", func(t *testing.T) {
		storage := inmemory_bookings_storage.New()
		competitorFlightsProvider := inmemory_competitor_flights_provider.New()
		launchpadRegistry := hardcoded_launchpad_registry.New()
		bookTicket := book_ticket.NewTicketBooker(storage, competitorFlightsProvider, launchpadRegistry)

		idFactory := id.NewChaoticFactory(t.Name())
		launchpadID := hardcoded_launchpad_registry.Texas
		destinationID := idFactory.Generate()
		launchDate := date.MustParse("25/11/2024")

		conflictingBooking := bookings.New((&id.FixedIDGenerator{}).Generate(),
			bookings.ClientData{
				FirstName: "John",
				LastName:  "Doe",
				Gender:    "Male",
				Birthday:  "07/03",
			}, destinationID.String(), launchpadID, launchDate)
		err := storage.SaveBooking(context.Background(), conflictingBooking)
		require.NoError(t, err)

		err = bookTicket.Execute(context.Background(),
			book_ticket.BookTicketParams{
				ID:            idFactory.Generate().String(),
				FirstName:     "Jane",
				LastName:      "Smith",
				Gender:        "Female",
				Birthday:      "29/09",
				LaunchpadID:   launchpadID,
				DestinationID: idFactory.Generate().String(), // other destination
				LaunchDate:    launchDate.String(),
			})

		assert.Error(t, err)
		assert.ErrorIs(t, err, book_ticket.ErrLaunchpadUnavailable)

		storedBookings := storage.ListBookings()
		require.GreaterOrEqual(t, len(storedBookings), 1)
		assert.Len(t, storedBookings, 1)

		firstBooking := storedBookings[0]
		assert.EqualValues(t, conflictingBooking.ToSnapshot(), firstBooking.ToSnapshot())
	})
}

func (s *testSuite) TestConflictWithCompetitor() {
	s.T().Run("Should not be able to book a ticket when the booking triggers a conflict with a competitor",
		func(t *testing.T) {

			competitorFlightsProvider := inmemory_competitor_flights_provider.New()
			bookingsStorage := inmemory_bookings_storage.New()
			launchpadRegistry := hardcoded_launchpad_registry.New()

			conflictingDate := date.MustParse("13/10/2024")
			err := competitorFlightsProvider.RegisterFlight(hardcoded_launchpad_registry.California,
				conflictingDate)
			require.NoError(t, err)

			bookTicket := book_ticket.NewTicketBooker(bookingsStorage, competitorFlightsProvider, launchpadRegistry)

			idFactory := id.NewChaoticFactory(t.Name())
			err = bookTicket.Execute(context.Background(),
				book_ticket.BookTicketParams{
					ID:            idFactory.Generate().String(),
					FirstName:     "John",
					LastName:      "Doe",
					Gender:        "Male",
					Birthday:      "12/05",
					LaunchpadID:   hardcoded_launchpad_registry.California,
					DestinationID: "example-destination-id",
					LaunchDate:    conflictingDate.String(),
				})

			assert.Error(t, err)
			assert.ErrorIs(t, err, book_ticket.ErrLaunchpadUnavailable)

			storedBookings := bookingsStorage.ListBookings()
			assert.Empty(t, storedBookings)
		})
}

func (s *testSuite) TestLaunchpadNotExisting() {
	s.T().Run("Should not be able to book a ticket when the launchpad does not exist",
		func(t *testing.T) {
			competitorFlightsProvider := inmemory_competitor_flights_provider.New()
			bookingsStorage := inmemory_bookings_storage.New()
			launchpadRegistry := hardcoded_launchpad_registry.New()
			bookTicket := book_ticket.NewTicketBooker(bookingsStorage, competitorFlightsProvider, launchpadRegistry)

			idFactory := id.NewChaoticFactory(t.Name())

			err := bookTicket.Execute(context.Background(),
				book_ticket.BookTicketParams{
					ID:            idFactory.Generate().String(),
					FirstName:     "John",
					LastName:      "Doe",
					Gender:        "Male",
					Birthday:      "12/05",
					LaunchpadID:   "does-not-exist",
					DestinationID: "example-destination-id",
					LaunchDate:    date.MustParse("25/11/2024").String(),
				})

			assert.Error(t, err)
			assert.ErrorIs(t, err, book_ticket.ErrLaunchpadDoesNotExist)

			storedBookings := bookingsStorage.ListBookings()
			assert.Empty(t, storedBookings)

		})
}
