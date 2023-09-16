package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseCoinSpent(data map[string]any, height pkgTypes.Level) (*storage.Address, error) {
	balance, err := getBalanceFromMap(data, "amount")
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, nil
	}

	if senderString := getStringFromMap(data, "spender"); senderString != "" {
		_, hash, err := pkgTypes.Address(senderString).Decode()
		if err != nil {
			return nil, errors.Wrapf(err, "decode sender: %s", senderString)
		}
		return &storage.Address{
			Address:    senderString,
			Hash:       hash,
			Height:     height,
			LastHeight: height,
			Balance: storage.Balance{
				Currency: balance.Denom,
				Total:    decimal.NewFromBigInt(balance.Amount.Neg().BigInt(), 0),
			},
		}, nil
	}

	return nil, nil
}

func parseCoinReceived(data map[string]any, height pkgTypes.Level) (*storage.Address, error) {
	balance, err := getBalanceFromMap(data, "amount")
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, nil
	}

	if receiverString := getStringFromMap(data, "receiver"); receiverString != "" {
		_, hash, err := pkgTypes.Address(receiverString).Decode()
		if err != nil {
			return nil, errors.Wrapf(err, "decode receiver: %s", receiverString)
		}
		return &storage.Address{
			Address:    receiverString,
			Hash:       hash,
			Height:     height,
			LastHeight: height,
			Balance: storage.Balance{
				Currency: balance.Denom,
				Total:    decimal.NewFromBigInt(balance.Amount.BigInt(), 0),
			},
		}, nil
	}

	return nil, nil
}
