package rollback

import (
	"context"
	"database/sql"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	indexerCfg "github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node/mock"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/status-im/keycard-go/hexutils"
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

func (s *ModuleTestSuite) TestModule_SuccessOnRollbackOneBlocks() {
	s.InitDb("../../../test/data/rollback")

	expectedHash := hexutils.HexToBytes("5F7A8DDFE6136FE76B65B9066D4F816D707F28C05B3362D66084664C5B39BA98")
	s.InitApi(func() {
		s.api.EXPECT().
			Block(gomock.Any(), types.Level(1000)).
			Return(types.ResultBlock{
				BlockID: tmTypes.BlockID{
					Hash: bytes.HexBytes{1}, // not equal with block in storage
				},
			}, nil).
			MaxTimes(1)

		s.api.EXPECT().
			Block(gomock.Any(), types.Level(999)).
			Return(types.ResultBlock{
				BlockID: tmTypes.BlockID{
					Hash: expectedHash,
				},
			}, nil).
			MaxTimes(1)
	})

	rollbackModule := NewModule(
		s.storage.Transactable,
		s.storage.State,
		s.storage.Blocks,
		s.api,
		indexerCfg.Indexer{Name: testIndexerName},
	)

	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stateListener := modules.New("state-listener")
	stateListener.CreateInput("state")
	err := stateListener.AttachTo(&rollbackModule, OutputName, "state")
	s.Require().NoError(err)

	rollbackModule.Start(ctx)
	defer func() {
		s.Require().NoError(rollbackModule.Close())
	}()

	// Act
	rollbackModule.MustInput(InputName).Push(struct{}{})

	for {
		select {
		case <-ctx.Done():
			s.T().Error("stop by cancelled context")
			return
		case msg, ok := <-stateListener.MustInput("state").Listen():
			s.Require().True(ok, "received value should be delivered by successful send operation")

			state, ok := msg.(storage.State)
			s.Require().True(ok, "got wrong type %T", msg)

			s.Require().Equal(types.Level(999), state.LastHeight)
			s.Require().Equal(expectedHash, state.LastHash)
			return
		}
	}
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
