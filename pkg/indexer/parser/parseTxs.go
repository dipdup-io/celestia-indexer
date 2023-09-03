package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

func parseTxs(b types.BlockData) ([]storage.Tx, error) {
	txs := make([]storage.Tx, len(b.TxsResults))

	for i, txRes := range b.TxsResults {
		t, err := parseTx(b, i, txRes)

		if err != nil {
			return nil, err
		}

		txs[i] = t
	}

	return txs, nil
}

func parseTx(b types.BlockData, index int, txRes *nodeTypes.ResponseDeliverTx) (storage.Tx, error) {
	d, err := decodeTx(b, index)
	if err != nil {
		return storage.Tx{}, errors.Wrapf(err, "parsing Tx on index %d error", index)
	}

	t := storage.Tx{
		Height:        b.Height,
		Time:          b.Block.Time,
		Position:      uint64(index),
		GasWanted:     uint64(txRes.GasWanted),
		GasUsed:       uint64(txRes.GasUsed),
		TimeoutHeight: d.timeoutHeight,
		EventsCount:   uint64(len(txRes.Events)),
		MessagesCount: uint64(len(d.messages)),
		Fee:           d.fee,
		Status:        storageTypes.StatusSuccess,
		Codespace:     txRes.Codespace,
		Hash:          b.Block.Txs[index].Hash(),
		Memo:          d.memo,
		MessageTypes:  storageTypes.MsgTypeBits{}, // TODO

		Messages: make([]storage.Message, len(d.messages)),
		Events:   nil,
	}

	if txRes.Code != 0 {
		t.Status = storageTypes.StatusFailed
		t.Error = txRes.Log
	}

	t.Events = parseEvents(b, txRes.Events)

	for position, msg := range d.messages {
		// TODO get namespace

		t.Messages[position] = storage.Message{
			Height:   b.Height,
			Time:     b.Block.Time,
			Position: uint64(position),
			Type:     storageTypes.MsgTypeUnknown, // TODO get type
			Data:     structs.Map(msg),
		}
	}

	return t, nil
}
