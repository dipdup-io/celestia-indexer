package receiver

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

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
