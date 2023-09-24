package receiver

import (
	"context"
	"database/sql"
	ic "github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node/mock"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/go-testfixtures/testfixtures/v3"
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
	api           *mock.MockApi
}

// SetupSuite -
func (s *ModuleTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 60*time.Second)
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
	s.api = mock.NewMockApi(ctrl)

	if configureApi != nil {
		configureApi()
	}
}

// func getResultBlock(hash types.Hex) types.ResultBlock {
// 	return types.ResultBlock{
// 		BlockID: types.BlockId{
// 			Hash: hash,
// 		},
// 	}
// }

var cfg = ic.Indexer{
	Name:         testIndexerName,
	ThreadsCount: 1,
	StartLevel:   0,
	BlockPeriod:  1,
}

func (s *ModuleTestSuite) createModule() Module {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	state, err := s.storage.State.ByName(ctx, testIndexerName)
	s.Require().NoError(err)

	receiverModule := NewModule(cfg, s.api, &state)

	return receiverModule
}

func (s *ModuleTestSuite) createModuleEmptyState() Module {
	receiverModule := NewModule(cfg, s.api, nil)
	return receiverModule
}

func (s *ModuleTestSuite) TestModule_SuccessOnStop() {
	s.InitDb("../../../test/data")
	s.InitApi(func() {
		s.api.EXPECT().Status(gomock.Any()).Return(nodeTypes.Status{}, nil).MinTimes(0)
	})

	receiverModule := s.createModule()

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	stopperModule := stopper.NewModule(cancelCtx)
	err := stopperModule.AttachTo(&receiverModule, StopOutput, stopper.InputName)
	s.Require().NoError(err)

	stopperCtx, stopperCtxCancel := context.WithCancel(context.Background())
	defer stopperCtxCancel()

	stopperModule.Start(stopperCtx)
	receiverModule.Start(ctx)

	defer func() {
		s.Require().NoError(receiverModule.Close())
	}()

	receiverModule.MustOutput(StopOutput).Push(struct{}{})

	for range ctx.Done() {
		s.Require().ErrorIs(context.Canceled, ctx.Err())
		return
	}

}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
