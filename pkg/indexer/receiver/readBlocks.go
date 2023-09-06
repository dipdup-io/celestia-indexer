package receiver

import (
	"bytes"
	"context"
	"encoding/hex"
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

	if rollbackDetected, err := r.checkRollback(ctx, headLevel, headHash); err != nil {
		return errors.Wrap(err, "while detecting rollback")
	} else if rollbackDetected {
		r.log.Debug().Msg("rollback detected")
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

func (r *Receiver) checkRollback(ctx context.Context, headLevel storage.Level, headHash []byte) (rollbackDetected bool, err error) {
	receiverLevel, receiverHash := r.Level()

	if receiverLevel == headLevel && !bytes.Equal(receiverHash, headHash) { // TODO-DISCUSS seems like it won't happen
		r.log.Info().
			Str("receiverHash", hex.EncodeToString(receiverHash)).
			Str("headHash", hex.EncodeToString(headHash)).
			Msg("rollback detected")

		rollbackDetected = true
		return
	}

	if receiverLevel != headLevel {
		block, e := r.api.Block(ctx, receiverLevel)
		if e != nil {
			err = errors.Wrapf(e, "while getting block by receiver level=%d", receiverLevel)
			return
		}

		if !bytes.Equal(receiverHash, block.BlockID.Hash) {
			r.log.Info().
				Str("receiverHash", hex.EncodeToString(receiverHash)).
				Str("blockHash", hex.EncodeToString(block.BlockID.Hash)).
				Msg("rollback detected")

			rollbackDetected = true
			return
		}
	}

	return
}
