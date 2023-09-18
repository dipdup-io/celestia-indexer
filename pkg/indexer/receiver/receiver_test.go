package receiver

import (
	"context"
	"database/sql"
	"github.com/dipdup-io/celestia-indexer/pkg/node/mock"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/tendermint/tendermint/libs/bytes"
	tmTypes "github.com/tendermint/tendermint/types"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/stretchr/testify/suite"
)

const testIndexerName = "test_indexer"

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       postgres.Storage
	api           *mock.MockAPI
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

	st, err := postgres.Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)
	s.storage = st
}

// TearDownSuite -
func (s *ModuleTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *ModuleTestSuite) InitDb(path string) {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory(path),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

func (s *ModuleTestSuite) InitApi(configureApi func()) {
	ctrl := gomock.NewController(s.T())
	s.api = mock.NewMockAPI(ctrl)

	if configureApi != nil {
		configureApi()
	}
}

func GetResultBlock(hash bytes.HexBytes) types.ResultBlock {
	return types.ResultBlock{
		BlockID: tmTypes.BlockID{
			Hash: hash,
		},
	}
}

func (s *ModuleTestSuite) TestModule_OnClosedInput() {
	s.InitDb("../../../test/data")
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
