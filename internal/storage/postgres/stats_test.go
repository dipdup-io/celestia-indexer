package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

// TODO: write test on all entities

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *StatsTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb:latest-pg15",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	strg, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)
	s.storage = strg

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory("../../../test/data"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *StatsTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *StatsTestSuite) TestCountBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Stats.Count(ctx, storage.CountRequest{
		Table: "block",
		From:  1672573739,
	})
	s.Require().NoError(err)
	s.Require().EqualValues("2", count)
}

func (s *StatsTestSuite) TestCountBlockNoData() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Stats.Count(ctx, storage.CountRequest{
		Table: "block",
		From:  1693324139,
	})
	s.Require().NoError(err)
	s.Require().EqualValues("0", count)
}

func (s *StatsTestSuite) TestSummaryBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	totalFee, err := s.storage.Stats.Summary(ctx, storage.SummaryRequest{
		CountRequest: storage.CountRequest{
			Table: "block",
			From:  1672573739,
		},
		Function: "sum",
		Column:   "fee",
	})
	s.Require().NoError(err)
	s.Require().Equal("4599819996", totalFee)
}

func (s *StatsTestSuite) TestHistogramBlock() {
	type test struct {
		timeframe storage.Timeframe
		wantDate  time.Time
	}

	tests := []test{
		{
			timeframe: storage.TimeframeHour,
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeDay,
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeWeek,
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeMonth,
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeYear,
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer ctxCancel()

		histogram, err := s.storage.Stats.Histogram(ctx,
			storage.HistogramRequest{
				SummaryRequest: storage.SummaryRequest{
					CountRequest: storage.CountRequest{
						Table: "block",
						From:  1672573739,
					},
					Function: "sum",
					Column:   "fee",
				},
				Timeframe: tests[i].timeframe,
			})
		s.Require().NoError(err)
		s.Require().Len(histogram, 1)

		item := histogram[0]
		s.Require().Equal("4599819996", item.Value)
		s.Require().True(item.Time.Equal(tests[i].wantDate))
	}
}

func (s *StatsTestSuite) TestHistogramCountBlock() {
	type test struct {
		timeframe storage.Timeframe
		wantDate  time.Time
	}

	tests := []test{
		{
			timeframe: storage.TimeframeHour,
			wantDate:  time.Date(2023, 7, 4, 3, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeDay,
			wantDate:  time.Date(2023, 7, 4, 0, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeWeek,
			wantDate:  time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeMonth,
			wantDate:  time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
		}, {
			timeframe: storage.TimeframeYear,
			wantDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for i := range tests {
		ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer ctxCancel()

		histogram, err := s.storage.Stats.HistogramCount(ctx,
			storage.HistogramCountRequest{
				CountRequest: storage.CountRequest{
					Table: "block",
					From:  1672573739,
				},
				Timeframe: tests[i].timeframe,
			})
		s.Require().NoError(err)
		s.Require().Len(histogram, 1)

		item := histogram[0]
		s.Require().Equal("2", item.Value)
		s.Require().True(item.Time.Equal(tests[i].wantDate))
	}
}

func TestSuiteStats_Run(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}
