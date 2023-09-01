package parser

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

func (p *Parser) parse(ctx context.Context, b types.BlockData) error {
	p.log.Info().Int64("height", b.Block.Height).Msg("parsing block...")

	block := storage.Block{
		Height:       b.Height,
		Time:         b.Block.Time,
		VersionBlock: b.Block.Version.Block,
		VersionApp:   b.Block.Version.App,

		TxCount:      uint64(len(b.Block.Data.Txs)),
		EventsCount:  uint64(len(b.BeginBlockEvents) + len(b.EndBlockEvents)),
		MessageTypes: storageTypes.MsgTypeBits{}, // TODO init
		BlobsSize:    0,

		Hash:               []byte(b.BlockID.Hash), // TODO create a Hex type for common usage through indexer app
		ParentHash:         []byte(b.Block.LastBlockID.Hash),
		LastCommitHash:     []byte(b.Block.LastCommitHash),
		DataHash:           []byte(b.Block.DataHash),
		ValidatorsHash:     []byte(b.Block.ValidatorsHash),
		NextValidatorsHash: []byte(b.Block.NextValidatorsHash),
		ConsensusHash:      []byte(b.Block.ConsensusHash),
		AppHash:            []byte(b.Block.AppHash),
		LastResultsHash:    []byte(b.Block.LastResultsHash),
		EvidenceHash:       []byte(b.Block.EvidenceHash),
		ProposerAddress:    []byte(b.Block.ProposerAddress),

		Fee:     decimal.Zero, // TODO sum of auth_info.fee // RESEARCH
		ChainId: b.Block.ChainID,

		Txs:    parseTxs(b),
		Events: nil, // TODO
	}

	p.output.Push(block)
	return nil
}
