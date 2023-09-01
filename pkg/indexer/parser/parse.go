package parser

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"strconv"
)

func (p *Parser) parse(ctx context.Context, b types.BlockData) error {
	p.log.Info().Int64("height", b.Block.Height).Msg("parsing block...")

	block := storage.Block{
		Height:       b.Height,
		Time:         b.Block.Time,
		VersionBlock: strconv.FormatUint(b.Block.Version.Block, 10), // should we use string in storage type?
		VersionApp:   strconv.FormatUint(b.Block.Version.App, 10),   // should we use string in storage type?

		TxCount:       uint64(len(b.Block.Data.Txs)),
		EventsCount:   0, // TODO
		MessageTypes:  storageTypes.MsgTypeBits{},
		NamespaceSize: 0, // "Summary block namespace size from pay for blob"` // should it be in block?

		Hash:               []byte(b.BlockID.Hash), // create a Hex type for common usage through indexer app
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

		Fee:     decimal.Zero, // TODO
		ChainId: b.Block.ChainID,

		Txs:    make([]storage.Tx, 0),    // TODO
		Events: make([]storage.Event, 0), // TODO
	}

	p.output.Push(block)
	return nil
}
