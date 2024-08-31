package psql_bookings_storage_test

import (
	"testing"

	"github.com/raphoester/space-trouble-api/internal/pkg/test_envs"
	"github.com/stretchr/testify/suite"
)

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite
	pg *test_envs.Postgres
}

func (s *testSuite) SetupSuite() {
	var err error
	s.pg, err = test_envs.NewPostgres()
	s.Require().NoError(err)
}

func (s *testSuite) SetupTest() {
	err := s.pg.Clean()
	s.Assert().NoError(err)
}

func (s *testSuite) TearDownSuite() {
	_ = s.pg.Destroy()
}

func (s *testSuite) TestExample() {

}
