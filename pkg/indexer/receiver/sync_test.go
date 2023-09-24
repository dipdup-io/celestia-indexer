package receiver

import (
	"context"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
)

func (s *ModuleTestSuite) TestModule_SyncGracefullyStops() {
	s.InitDb("../../../test/data/empty")
	s.InitApi(func() {
		s.api.EXPECT().
			Status(gomock.Any()).
			Return(nodeTypes.Status{}, errors.New("service is down")).
			MaxTimes(1)
	})

	receiverModule := s.createModule()

	// ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	stopperModule := stopper.NewModule(cancelCtx)
	err := stopperModule.AttachTo(&receiverModule, StopOutput, stopper.InputName)
	s.Require().NoError(err)

	stopperCtx, stopperCtxCancel := context.WithCancel(context.Background())
	defer stopperCtxCancel()

	stopperModule.Start(stopperCtx)

	workersCtx, cancelWorkers := context.WithCancel(ctx)
	receiverModule.cancelWorkers = cancelWorkers
	receiverModule.pool.Start(workersCtx)

	go receiverModule.sync(ctx)

	defer close(receiverModule.blocks)

	for range ctx.Done() {
		s.Require().ErrorIs(context.Canceled, ctx.Err())
		return
	}
}
