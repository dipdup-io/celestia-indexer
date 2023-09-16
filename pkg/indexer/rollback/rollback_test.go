package rollback

import (
	"context"
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/stretchr/testify/suite"
)

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       postgres.Storage
}

// SetupSuite -
func (s *ModuleTestSuite) SetupSuite() {
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

	storage, err := postgres.Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)
	s.storage = storage
}

// TearDownSuite -
func (s *ModuleTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *ModuleTestSuite) TestModule_SuccessOnRollbackTwoBlocks() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory("../../../test/data/rollback"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}