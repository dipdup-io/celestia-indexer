package receiver

import (
	"context"
	"github.com/pkg/errors"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func (r *Receiver) readBlocks(ctx context.Context) error {
	headLevel, err := r.headLevel(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	for ; r.level <= headLevel; r.level++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			r.pool.AddTask(r.level)
		}
	}

	return nil
}

func (r *Receiver) headLevel(ctx context.Context) (storage.Level, error) {
	head, err := r.api.Head(ctx) // TODO read from status, get hash also
	if err != nil {
		return 0, err
	}

	headLevel := storage.Level(head.Block.Height)
	return headLevel, nil
}
