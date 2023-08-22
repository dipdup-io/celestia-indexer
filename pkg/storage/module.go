package storage

import (
	"context"
	"sync"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InputName -
const InputName = "data"

const defaultIndexerName = "celestia-indexer"

// Module -
type Module struct {
	storage postgres.Storage
	input   *modules.Input
	state   *storage.State
	log     zerolog.Logger
	wg      *sync.WaitGroup
}

// NewModule -
func NewModule(pg postgres.Storage, opts ...ModuleOption) Module {
	m := Module{
		storage: pg,
		input:   modules.NewInput(InputName),
		state: &storage.State{
			Name: defaultIndexerName,
		},
		wg: new(sync.WaitGroup),
	}
	m.log = log.With().Str("module", m.Name()).Logger()

	for i := range opts {
		opts[i](&m)
	}

	return m
}

// Name -
func (Module) Name() string {
	return "storage"
}

func (module *Module) initState(ctx context.Context) error {
	module.log.Info().Msg("loading current state from database...")

	state, err := module.storage.State.ByName(ctx, module.state.Name)
	switch {
	case err == nil:
		module.state = &state
		module.log.Info().
			Str("indexer_name", module.state.Name).
			Uint64("height", module.state.LastHeight).
			Time("last_updated", module.state.LastTime).
			Msg("current state")
		return nil

	case module.storage.State.IsNoRows(err):
		module.log.Info().Msg("state is not found. creating empty state...")
		return module.storage.State.Update(ctx, module.state)

	default:
		return errors.Wrap(err, "state loading")
	}
}

// Start -
func (module *Module) Start(ctx context.Context) {
	if err := module.initState(ctx); err != nil {
		module.log.Err(err).Msg("error during storage module initialization")
		return
	}

	module.wg.Add(1)
	go module.listen(ctx)
}

func (module *Module) listen(ctx context.Context) {
	defer module.wg.Done()

	module.log.Info().Msg("module started")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-module.input.Listen():
			if !ok {
				module.log.Warn().Msg("can't read message from input")
				continue
			}
			block, ok := msg.(storage.Block)
			if !ok {
				module.log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			if err := module.saveBlock(ctx, block); err != nil {
				module.log.Err(err).Msg("block saving error")
				continue
			}

			if err := module.notify(ctx, block); err != nil {
				module.log.Err(err).Msg("block notification error")
			}
		}
	}
}

// Close -
func (module *Module) Close() error {
	module.log.Info().Msg("closing module...")
	module.wg.Wait()

	return module.input.Close()
}

// Output -
func (module *Module) Output(name string) (*modules.Output, error) {
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
}

// Input -
func (module *Module) Input(name string) (*modules.Input, error) {
	if name != InputName {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return module.input, nil
}

// AttachTo -
func (module *Module) AttachTo(name string, input *modules.Input) error {
	output, err := module.Output(name)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}

func (module *Module) updateState(block storage.Block) {
	if block.Id <= module.state.LastHeight {
		return
	}

	module.state.LastHeight = block.Height
	module.state.LastTime = block.Time
	module.state.TotalTx += block.TxCount
	// TODO: update rest fields
}

func (module *Module) saveBlock(ctx context.Context, block storage.Block) error {
	module.log.Info().Uint64("height", block.Id).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage.Transactable)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := tx.Add(ctx, &block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveTransactions(ctx, block.Txs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var (
		messages = make([]any, 0)
		events   = make([]any, len(block.Events))
	)

	for i := range block.Events {
		events[i] = &block.Events[i]
	}

	for i := range block.Txs {
		for j := range block.Txs[i].Messages {
			block.Txs[i].Messages[j].TxId = block.Txs[i].Id
			messages = append(messages, &block.Txs[i].Messages[j])
		}

		for j := range block.Txs[i].Events {
			block.Txs[i].Events[j].TxId = &block.Txs[i].Id
			events = append(events, &block.Txs[i].Events[j])
		}
	}

	if len(messages) > 0 {
		if err := tx.BulkSave(ctx, messages); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if len(events) > 0 {
		if err := tx.BulkSave(ctx, events); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	// TODO: save addresses and namespaces

	module.updateState(block)
	if err := tx.Update(ctx, module.state); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	module.log.Info().Uint64("height", block.Id).Msg("block saved")
	return nil
}

func (module *Module) notify(ctx context.Context, block storage.Block) error {
	data, err := json.MarshalContext(ctx, block, json.UnorderedMap())
	if err != nil {
		return err
	}
	if err := module.storage.Notificator.Notify(ctx, storage.ChannelHead, string(data)); err != nil {
		return err
	}
	return nil
}
