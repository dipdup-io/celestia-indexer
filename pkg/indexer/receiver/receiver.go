package receiver

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

const (
	name         = "receiver"
	BlocksOutput = "blocks"
)

type Module struct {
	api     node.API
	cfg     config.Indexer
	outputs map[string]*modules.Output
	pool    *workerpool.Pool[storage.Level]
	blocks  chan types.BlockData
	level   storage.Level
	hash    []byte
	mx      *sync.RWMutex
	log     zerolog.Logger
	wg      *sync.WaitGroup
}

func NewModule(cfg config.Indexer, api node.API, state *storage.State) Module {
	var level storage.Level
	var hash []byte

	if state == nil {
		level = storage.Level(cfg.StartLevel)
		// TODO-DISCUSS check for hash changed of state last block
	} else {
		level = state.LastHeight
		hash = state.LastHash
	}

	receiver := Module{
		api:     api,
		cfg:     cfg,
		outputs: map[string]*modules.Output{BlocksOutput: modules.NewOutput(BlocksOutput)},
		blocks:  make(chan types.BlockData, cfg.ThreadsCount*10),
		level:   level,
		hash:    hash,
		mx:      new(sync.RWMutex),
		log:     log.With().Str("module", name).Logger(),
		wg:      new(sync.WaitGroup),
	}

	receiver.pool = workerpool.NewPool(receiver.worker, int(cfg.ThreadsCount))

	return receiver
}

// Name -
func (*Module) Name() string {
	return name
}

func (r *Module) Start(ctx context.Context) {
	r.log.Info().Msg("starting receiver...")
	r.pool.Start(ctx)

	r.wg.Add(1)
	go r.sequencer(ctx)

	r.wg.Add(1)
	go r.sync(ctx)
}

func (r *Module) Close() error {
	r.log.Info().Msg("closing...")
	r.wg.Wait()

	if err := r.pool.Close(); err != nil {
		return err
	}

	close(r.blocks)

	return nil
}

func (r *Module) Output(name string) (*modules.Output, error) {
	output, ok := r.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return output, nil
}

func (r *Module) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

func (r *Module) AttachTo(outputName string, input *modules.Input) error {
	output, err := r.Output(outputName)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}

func (r *Module) Level() (storage.Level, []byte) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.level, r.hash
}

func (r *Module) setLevel(level storage.Level, hash []byte) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.level = level
	r.hash = hash
}
