package receiver

import (
	"context"
	"github.com/pkg/errors"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func (r *Receiver) readBlocks(ctx context.Context) error {
	headLevel, headHash, err := r.head(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	for level := r.level; level <= headLevel; level++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			r.pool.AddTask(level)
		}
	}

	if rollbackDetected, err := r.checkRollback(ctx, headLevel, headHash); err != nil {
		return errors.Wrap(err, "while detecting rollback")
	} else if rollbackDetected {
		r.log.Debug().Msg("rollback detected")
		// TODO	call rollback to the rescue
	}

	return nil
}

func (r *Receiver) head(ctx context.Context) (storage.Level, []byte, error) {
	status, err := r.api.Status(ctx)
	if err != nil {
		return 0, nil, err
	}

	headLevel := storage.Level(status.SyncInfo.LatestBlockHeight)
	return headLevel, status.SyncInfo.LatestBlockHash, nil
}
