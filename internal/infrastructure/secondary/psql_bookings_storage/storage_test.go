package psql_bookings_storage_test

import (
	"context"
	"testing"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/psql_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/pkg/test_envs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite
	pg      *test_envs.Postgres
	storage *psql_bookings_storage.Storage
}

func (s *testSuite) SetupSuite() {
	var err error
	s.pg, err = test_envs.NewPostgres()
	s.Require().NoError(err)
	s.storage = psql_bookings_storage.New(s.pg.PG)
}

func (s *testSuite) SetupTest() {
	err := s.pg.Clean()
	s.Assert().NoError(err)
}

func (s *testSuite) TearDownSuite() {
	_ = s.pg.Destroy()
}

func exampleBookingSnapshot(opts ...func(snapshot *bookings.BookingSnapshot)) *bookings.BookingSnapshot {
	snapshot := &bookings.BookingSnapshot{
		ID:            "example-id",
		FirstName:     "John",
		LastName:      "Doe",
		Gender:        "Male",
		Birthday:      "19/10",
		DestinationID: "example-destination-id",
		LaunchpadID:   "example-launchpad-id",
		LaunchDate:    "10/10/2024",
	}

	for _, opt := range opts {
		opt(snapshot)
	}

	return snapshot
}

func (s *testSuite) TestSaveBooking() {
	s.T().Run("should save a booking", func(t *testing.T) {
		restored, err := bookings.Restore(exampleBookingSnapshot())
		require.NoError(t, err)

		err = s.storage.SaveBooking(context.Background(), restored)
		require.NoError(t, err)

		inDB, err := s.storage.Get(context.Background(), restored.ID())
		assert.NoError(t, err)
		assert.EqualValues(t, restored.ToSnapshot(), inDB.ToSnapshot())
	})
}

func (s *testSuite) TestListNoConflictNoOtherBooking() {
	s.T().Run("should not find any conflict since no other booking is on the table", func(t *testing.T) {
		restored, err := bookings.Restore(exampleBookingSnapshot())
		require.NoError(t, err)

		conflicting, err := s.storage.ListConflictingFlightBookings(context.Background(), restored)
		require.NoError(t, err)
		assert.Empty(t, conflicting)
	})
}

func (s *testSuite) TestListNoConflictDifferentDate() {
	s.T().Run("should not find a conflict since the launchpad is used at another date", func(t *testing.T) {
		restored, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
			snapshot.ID = "id1"
			snapshot.LaunchDate = "11/10/2024"
			snapshot.DestinationID = "destination1"
		}))
		require.NoError(t, err)

		err = s.storage.SaveBooking(context.Background(), restored)
		require.NoError(t, err)

		restored2, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
			snapshot.ID = "id2"
			snapshot.LaunchDate = "12/10/2024"
			snapshot.DestinationID = "destination2"
		}))
		require.NoError(t, err)

		conflicting, err := s.storage.ListConflictingFlightBookings(context.Background(), restored2)
		require.NoError(t, err)
		assert.Empty(t, conflicting)
	})
}

func (s *testSuite) TestListNoConflictDifferentLaunchpad() {
	s.T().Run("should not find a conflict since another launch is made from another launchpad",
		func(t *testing.T) {
			restored, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
				snapshot.ID = "id1"
				snapshot.LaunchpadID = "launchpad1"
				snapshot.DestinationID = "destination1"
			}))
			require.NoError(t, err)

			err = s.storage.SaveBooking(context.Background(), restored)
			require.NoError(t, err)

			restored2, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
				snapshot.ID = "id2"
				snapshot.LaunchpadID = "launchpad2"
				snapshot.DestinationID = "destination2"
			}))
			require.NoError(t, err)

			conflicting, err := s.storage.ListConflictingFlightBookings(context.Background(), restored2)
			require.NoError(t, err)
			assert.Empty(t, conflicting)
		})
}

func (s *testSuite) TestNoConflictSameDestination() {
	s.T().Run("should not find a conflict since the destination is the same between both bookings' flights",
		func(t *testing.T) {
			restored, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
				snapshot.ID = "id1"
				snapshot.DestinationID = "destination1"
			}))
			require.NoError(t, err)

			err = s.storage.SaveBooking(context.Background(), restored)
			require.NoError(t, err)

			restored2, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
				snapshot.ID = "id2"
				snapshot.DestinationID = "destination1"
			}))
			require.NoError(t, err)

			conflicting, err := s.storage.ListConflictingFlightBookings(context.Background(), restored2)
			require.NoError(t, err)
			assert.Empty(t, conflicting)
		})
}

func (s *testSuite) TestListConflict() {
	s.T().Run("should find a conflict since the launchpad is used for another destination at the same date",
		func(t *testing.T) {
			restored, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
				snapshot.ID = "id1"
				snapshot.LaunchDate = "10/10/2024"
				snapshot.DestinationID = "destination1"
			}))
			require.NoError(t, err)

			err = s.storage.SaveBooking(context.Background(), restored)
			require.NoError(t, err)

			restored2, err := bookings.Restore(exampleBookingSnapshot(func(snapshot *bookings.BookingSnapshot) {
				snapshot.ID = "id2"
				snapshot.LaunchDate = "10/10/2024"
				snapshot.DestinationID = "destination2"
			}))
			require.NoError(t, err)

			conflicting, err := s.storage.ListConflictingFlightBookings(context.Background(), restored2)
			require.NoError(t, err)
			require.Len(t, conflicting, 1)
			assert.EqualValues(t, restored.ToSnapshot(), conflicting[0].ToSnapshot())
		},
	)
}
