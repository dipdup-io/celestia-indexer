package parser

import (
	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	tmTypes "github.com/tendermint/tendermint/types"
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
	_, timeoutHeight, memo, msgs, err := decodeTx(b, index) // authInfo
	if err != nil {
		return storage.Tx{}, errors.Wrapf(err, "parsing Tx on index %d error", index)
	}

	t := storage.Tx{
		Height:        b.Height,
		Time:          b.Block.Time,
		Position:      uint64(index),
		GasWanted:     uint64(txRes.GasWanted),
		GasUsed:       uint64(txRes.GasUsed),
		TimeoutHeight: timeoutHeight,
		EventsCount:   uint64(len(txRes.Events)),
		MessagesCount: uint64(len(msgs)),
		Fee:           decimal.Zero, // TODO like nodes.guru
		Status:        storageTypes.StatusSuccess,
		Codespace:     txRes.Codespace,
		Hash:          b.Block.Txs[index].Hash(),
		Memo:          memo,                       // TODO like nodes.guru
		MessageTypes:  storageTypes.MsgTypeBits{}, // TODO

		Messages: nil, // make([]storage.Message, 0), // TODO
		Events:   nil,
	}

	if txRes.Code != 0 {
		t.Status = storageTypes.StatusFailed
		t.Error = txRes.Log
	}

	t.Events = parseEvents(b, txRes.Events)

	return t, nil
}

func decodeTx(b types.BlockData, index int) (authInfo tx.AuthInfo, timeoutHeight uint64, memo string, msgs []cosmosTypes.Msg, err error) {
	txCfg := createDecoder()
	decoder := txCfg.TxConfig.TxDecoder()

	raw := b.Block.Txs[index]
	if bTx, isBlob := tmTypes.UnmarshalBlobTx(raw); isBlob {
		raw = bTx.Tx
	}

	var txRaw txTypes.TxRaw
	if e := txCfg.Codec.Unmarshal(raw, &txRaw); e != nil {
		err = errors.Wrap(e, "unmarshaling tx error")
		return
	}
	if e := txCfg.Codec.Unmarshal(txRaw.AuthInfoBytes, &authInfo); e != nil {
		err = errors.Wrap(e, "decoding tx auth_info error")
		return
	}

	txDecoded, e := decoder(raw)
	if e != nil {
		err = errors.Wrap(e, "decoding tx error")
		return
	}

	if t, ok := txDecoded.(cosmosTypes.TxWithTimeoutHeight); ok {
		timeoutHeight = t.GetTimeoutHeight()
	}
	if t, ok := txDecoded.(cosmosTypes.TxWithMemo); ok {
		memo = t.GetMemo()
	}

	msgs = txDecoded.GetMsgs()

	// for _, msg := range txDecoded.GetMsgs() {
	//}

	return
}

func createDecoder() encoding.Config {
	return encoding.MakeConfig(app.ModuleEncodingRegisters...)
}
