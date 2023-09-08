package receiver

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func (r *Module) sequencer(ctx context.Context) {
	orderedBlocks := map[int64]types.BlockData{}
	var prevBlockHash []byte
	l, _ := r.Level()
	currentBlock := int64(l)

	for {
		select {
		case <-ctx.Done():
			return
		case block := <-r.blocks:
			orderedBlocks[block.Block.Height] = block

			if currentBlock == 0 {
				if err := r.receiveGenesis(ctx); err != nil {
					return
					// TODO: handle error on getting genesis, stop indexer
				}

				currentBlock += 1
				break
			}

			if b, ok := orderedBlocks[currentBlock]; ok {
				if prevBlockHash != nil {
					if !bytes.Equal(b.Block.LastBlockID.Hash, prevBlockHash) {
						prevBlockHash, currentBlock, orderedBlocks = r.startRollback(ctx, b, prevBlockHash)
						break
					}
				} // TODO else: check with block from storage?

				r.outputs[BlocksOutput].Push(b)
				r.setLevel(types.Level(currentBlock), b.BlockID.Hash)
				r.log.Debug().Msgf("put in order block=%d", currentBlock)

				prevBlockHash = b.BlockID.Hash
				delete(orderedBlocks, currentBlock)
				currentBlock += 1
			}
		}
	}
}

func (r *Module) startRollback(
	ctx context.Context,
	b types.BlockData,
	prevBlockHash []byte,
) ([]byte, int64, map[int64]types.BlockData) {
	r.log.Info().
		Str("current.lastBlockHash", hex.EncodeToString(b.Block.LastBlockID.Hash)).
		Str("prevBlockHash", hex.EncodeToString(prevBlockHash)).
		Uint64("level", uint64(b.Height)).
		Msg("rollback detected")

	r.rollbackSync.Add(1)
	if r.cancelReadBlocks != nil {
		r.cancelReadBlocks()
	}
	r.outputs[RollbackOutput].Push(struct{}{})
	r.rollbackSync.Wait()

	l, _ := r.Level()
	currentBlock := int64(l)
	prevBlockHash = nil
	orderedBlocks := map[int64]types.BlockData{}

	clearChannel(r.blocks)
	workersCtx, cancelWorkers := context.WithCancel(ctx)
	r.cancelWorkers = cancelWorkers
	r.pool.Start(workersCtx)

	return prevBlockHash, currentBlock, orderedBlocks
}

func clearChannel(blocks <-chan types.BlockData) {
	for len(blocks) > 0 {
		<-blocks
	}
}
