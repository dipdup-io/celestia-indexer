package receiver

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/workerpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Receiver struct {
	// api
	timeout time.Duration
	// output       *modules.Output
	pool         *workerpool.Pool[storage.Level]
	processing   map[storage.Level]struct{}
	processingMx *sync.Mutex
	log          zerolog.Logger
	wg           *sync.WaitGroup
}

func New(cfg config.Config) *Receiver {
	receiver := &Receiver{
		processing:   make(map[storage.Level]struct{}),
		processingMx: new(sync.Mutex),
		log:          log.With().Str("module", "receiver").Logger(),
		timeout:      time.Duration(cfg.Indexer.Timeout) * time.Second,
		wg:           new(sync.WaitGroup),
	}

	if receiver.timeout == 0 {
		receiver.timeout = 10 * time.Second
	}

	receiver.pool = workerpool.NewPool(receiver.worker, int(cfg.Indexer.ThreadsCount))

	return receiver
}

func (r *Receiver) worker(ctx context.Context, level storage.Level) {

}
