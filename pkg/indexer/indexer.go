package indexer

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/parser"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/node/rpc"
	"github.com/dipdup-io/celestia-indexer/pkg/storage"
	"sync"

	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/receiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg      config.Config
	api      node.API
	receiver *receiver.Receiver
	parser   *parser.Parser
	storage  *storage.Module
	wg       *sync.WaitGroup
	log      zerolog.Logger
}

func New(ctx context.Context, cfg config.Config) *Indexer {

	api := rpc.NewAPI(cfg.DataSources["node_rpc"])
	r := receiver.NewModule(cfg, &api)

	p := parser.NewModule()

	pg, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		log.Err(err).Msg("creating pg context in indexer")
	}
	s := storage.NewModule(pg)

	return &Indexer{
		cfg:      cfg,
		api:      &api,
		receiver: &r,
		parser:   &p,
		storage:  &s,
		wg:       new(sync.WaitGroup),
		log:      log.With().Str("module", "indexer").Logger(),
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	i.log.Info().Msg("starting indexer...")

	i.receiver.Start(ctx)

	return nil
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("closing indexer...")
	i.wg.Wait()

	if err := i.receiver.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}

	return nil
}
