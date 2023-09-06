package receiver

import (
	"bytes"
	"context"
	"github.com/pkg/errors"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

func (r *Receiver) readBlocks(ctx context.Context) error {
	headLevel, headHash, err := r.headLevel(ctx)
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

	receiverLevel, receiverHash := r.Level()
	if receiverLevel == headLevel && !bytes.Equal(receiverHash, headHash) {
		r.log.Info().
			Bytes("receiverHash", receiverHash).
			Bytes("headHash", headHash).
			Msg("rollback detected")
		// TODO	call rollback to the rescue
	}

	return nil
}

func (r *Receiver) headLevel(ctx context.Context) (storage.Level, []byte, error) {
	head, err := r.api.Head(ctx) // TODO read from status, get hash also
	if err != nil {
		return 0, nil, err
	}

	headLevel := storage.Level(head.Block.Height)
	return headLevel, head.BlockID.Hash, nil
}
