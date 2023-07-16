//go:build integration

package integrationtest

import (
	"context"
	"route256/checkout/internal/config"
	"route256/checkout/internal/repository/postgres"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

// Group integration tests and data for it
type Suite struct {
	suite.Suite
	pg   *pgxpool.Pool
	cart *postgres.CartRepository
}

// Starting point for tests
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// Setup environment for integration tests
func (s *Suite) SetupSuite() {
	cfg, err := config.New()
	s.Require().NoError(err)

	s.pg, err = pgxpool.Connect(context.Background(), cfg.Postgres.TestDBConnectionString)
	s.Require().NoError(err)

	s.cart = postgres.New(s.pg)
}

// Clean db tables before each test
func (s *Suite) SetupTest() {
	query := "TRUNCATE TABLE "
	_, err := s.pg.Exec(context.Background(), query+tableNameCartItem)
	s.Require().NoError(err)
	_, err = s.pg.Exec(context.Background(), query+tableNameCart)
	s.Require().NoError(err)
}

// Tear down environment for integration tests after all tests
func (s *Suite) TearDownSuite() {
	query := "TRUNCATE TABLE "
	_, err := s.pg.Exec(context.Background(), query+tableNameCartItem)
	s.Require().NoError(err)
	_, err = s.pg.Exec(context.Background(), query+tableNameCart)
	s.Require().NoError(err)
	s.pg.Close()
}
