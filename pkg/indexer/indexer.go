package indexer

import (
	"context"
	"sync"

	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/receiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg      config.Config
	receiver *receiver.Receiver
	wg       *sync.WaitGroup
	log      zerolog.Logger
}

func New(cfg config.Config) *Indexer {

	return &Indexer{
		cfg:      cfg,
		receiver: receiver.New(cfg),
		wg:       new(sync.WaitGroup),
		log:      log.With().Str("module", "indexer").Logger(),
	}
}

func (i *Indexer) Start(ctx context.Context) error {
	i.log.Info().Msg("starting indexer...")
	return nil
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("stopping indexer...")
	i.wg.Wait()

	return nil
}
