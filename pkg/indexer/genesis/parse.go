package genesis

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type parsedData struct {
	block     storage.Block
	addresses map[string]*storage.Address
}

func newParsedData() parsedData {
	return parsedData{
		addresses: make(map[string]*storage.Address),
	}
}

func (module *Module) parse(genesis types.Genesis) (parsedData, error) {
	data := newParsedData()
	block := storage.Block{
		Time:    genesis.GenesisTime,
		Height:  pkgTypes.Level(genesis.InitialHeight - 1),
		AppHash: []byte(genesis.AppHash),
		ChainId: genesis.ChainID,
		Txs:     make([]storage.Tx, 0),
		Stats: storage.BlockStats{
			Time:          genesis.GenesisTime,
			Height:        pkgTypes.Level(genesis.InitialHeight - 1),
			TxCount:       uint64(len(genesis.AppState.Genutil.GenTxs)),
			EventsCount:   0,
			Fee:           decimal.Zero,
			SupplyChange:  decimal.Zero,
			InflationRate: decimal.Zero,
		},
	}

	for index, genTx := range genesis.AppState.Genutil.GenTxs {
		txDecoded, err := decode.JsonTx(genTx)
		if err != nil {
			return data, errors.Wrapf(err, "failed to decode GenTx '%s'", genTx)
		}

		memoTx, ok := txDecoded.(cosmosTypes.TxWithMemo)
		if !ok {
			return data, errors.Wrapf(err, "expected TxWithMemo, got %T", genTx)
		}
		txWithTimeoutHeight, ok := txDecoded.(cosmosTypes.TxWithTimeoutHeight)
		if !ok {
			return data, errors.Wrapf(err, "expected TxWithTimeoutHeight, got %T", genTx)
		}

		tx := storage.Tx{
			Height:        block.Height,
			Time:          block.Time,
			Position:      uint64(index),
			TimeoutHeight: txWithTimeoutHeight.GetTimeoutHeight(),
			MessagesCount: uint64(len(txDecoded.GetMsgs())),
			Fee:           decimal.Zero,
			Status:        storageTypes.StatusSuccess,
			Memo:          memoTx.GetMemo(),
			MessageTypes:  storageTypes.NewMsgTypeBitMask(),

			Messages:  make([]storage.Message, len(txDecoded.GetMsgs())),
			Events:    nil,
			Addresses: make([]storage.AddressWithType, 0),
		}

		for msgIndex, msg := range txDecoded.GetMsgs() {
			decoded, err := decode.Message(msg, block.Height, block.Time, msgIndex)
			if err != nil {
				return data, errors.Wrap(err, "decode genesis message")
			}

			tx.Messages[msgIndex] = decoded.Msg
			tx.MessageTypes.SetBit(decoded.Msg.Type)
			tx.BlobsSize += decoded.BlobsSize
			tx.Addresses = append(tx.Addresses, decoded.Addresses...)
		}

		block.Txs = append(block.Txs, tx)
	}

	module.parseTotalSupply(genesis.AppState.Bank.Supply, &block)

	if err := module.parseAccounts(genesis.AppState.Auth.Accounts, block.Height, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis accounts")
	}
	if err := module.parseBalances(genesis.AppState.Bank.Balances, block.Height, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis account balances")
	}

	data.block = block
	return data, nil
}

func (module *Module) parseTotalSupply(supply []types.Supply, block *storage.Block) {
	if len(supply) == 0 {
		return
	}

	if totalSupply, err := decimal.NewFromString(supply[0].Amount); err == nil {
		block.Stats.SupplyChange = totalSupply
	}
}

func (module *Module) parseAccounts(accounts []types.Accounts, height pkgTypes.Level, data *parsedData) error {
	for i := range accounts {
		address := storage.Address{
			Height:  height,
			Hash:    []byte(accounts[i].Address),
			Balance: decimal.Zero,
		}
		data.addresses[address.String()] = &address
	}
	return nil
}

func (module *Module) parseBalances(balances []types.Balances, height pkgTypes.Level, data *parsedData) error {
	for i := range balances {
		if len(balances[i].Coins) == 0 {
			continue
		}
		address := storage.Address{
			Hash:    []byte(balances[i].Address),
			Height:  height,
			Balance: decimal.Zero,
		}
		if balance, err := decimal.NewFromString(balances[i].Coins[0].Amount); err == nil {
			address.Balance = address.Balance.Add(balance)
		}

		if addr, ok := data.addresses[address.String()]; ok {
			addr.Balance = addr.Balance.Add(address.Balance)
		} else {
			data.addresses[address.String()] = &address
		}
	}

	return nil
}
