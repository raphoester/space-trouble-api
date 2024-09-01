package psql_bookings_getter_test

import (
	"context"
	"testing"

	"github.com/raphoester/space-trouble-api/internal/domain/model/bookings"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/psql_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/pkg/test_envs"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings/psql_bookings_getter"
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

func (s *testSuite) TestGetAllBookings() {
	s.T().Run("should find all bookings", func(t *testing.T) {
		snapshot := &bookings.BookingSnapshot{
			ID:            "example-id-1",
			FirstName:     "John",
			LastName:      "Doe",
			Gender:        "Male",
			Birthday:      "04/04",
			DestinationID: "mars",
			LaunchpadID:   "abc",
			LaunchDate:    "11/09/2024",
		}

		booking1, err := bookings.Restore(snapshot)
		require.NoError(t, err)

		snapshot.ID = "example-id-2"
		snapshot.LaunchpadID = "def"
		snapshot.DestinationID = "moon"
		snapshot.LaunchDate = "12/09/2024"

		err = s.storage.SaveBooking(context.Background(), booking1)
		require.NoError(t, err)

		booking2, err := bookings.Restore(snapshot)
		require.NoError(t, err)

		err = s.storage.SaveBooking(context.Background(), booking2)
		require.NoError(t, err)

		getter := psql_bookings_getter.New(s.pg.PG)
		retrievedBookings, err := getter.GetAllBookings(context.Background())
		require.NoError(t, err)

		require.Len(t, retrievedBookings, 2)

		for i := range retrievedBookings { // TODO: handle that through ordering in query
			if retrievedBookings[i].Id == "example-id-1" {
				assertDTOEqualSnapshot(t, retrievedBookings[i], booking1.ToSnapshot())
				continue
			}
			assertDTOEqualSnapshot(t, retrievedBookings[i], booking2.ToSnapshot())
		}
	})
}

func (s *testSuite) TestGetAllBookingsEmpty() {
	s.T().Run("should return empty list when no bookings", func(t *testing.T) {
		getter := psql_bookings_getter.New(s.pg.PG)
		retrievedBookings, err := getter.GetAllBookings(context.Background())
		require.NoError(t, err)
		require.Empty(t, retrievedBookings)
	})
}

func assertDTOEqualSnapshot(t *testing.T, dto get_all_bookings.Booking, snapshot bookings.BookingSnapshot) {
	t.Helper()
	require.Equal(t, snapshot.ID, dto.Id)
	require.Equal(t, snapshot.FirstName, dto.FirstName)
	require.Equal(t, snapshot.LastName, dto.LastName)
	require.Equal(t, snapshot.LaunchpadID, dto.LaunchpadID)
	require.Equal(t, snapshot.DestinationID, dto.DestinationID)
	require.Equal(t, snapshot.LaunchDate, dto.LaunchDate)
}
