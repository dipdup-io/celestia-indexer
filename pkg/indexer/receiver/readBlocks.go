package receiver

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func (r *Receiver) readBlocks(ctx context.Context) error {
	headLevel, err := r.headLevel(ctx)
	if err != nil {
		return err
	}

	startLevel := storage.Level(r.cfg.Indexer.StartLevel) // TODO read from current state
	level := startLevel

	for level <= headLevel {
		for ; level <= headLevel; level++ {
			select {
			case <-ctx.Done():
				return nil
			default:
				r.pool.AddTask(level)
			}
		}

		if headLevel, err = r.headLevel(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Receiver) headLevel(ctx context.Context) (storage.Level, error) {
	head, err := r.api.Head(ctx)
	if err != nil {
		return 0, err
	}

	headLevel := storage.Level(head.Block.Height)
	return headLevel, nil
}
