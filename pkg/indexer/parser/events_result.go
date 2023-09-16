package parser

import (
	"strings"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type eventsResult struct {
	SupplyChange  decimal.Decimal
	InflationRate decimal.Decimal

	Addresses []storage.Address
}

func (er *eventsResult) Fill(events []storage.Event) error {
	er.Addresses = make([]storage.Address, 0)

	for i := range events {
		switch events[i].Type {
		case types.EventTypeBurn:
			er.SupplyChange = er.SupplyChange.Sub(getDecimalFromMap(events[i].Data, "amount"))
		case types.EventTypeMint:
			er.InflationRate = getDecimalFromMap(events[i].Data, "inflation_rate")
			er.SupplyChange = er.SupplyChange.Add(getDecimalFromMap(events[i].Data, "amount"))
		case types.EventTypeCoinReceived:
			address, err := parseCoinReceived(events[i].Data, events[i].Height)
			if err != nil {
				return errors.Wrap(err, "parse coin received")
			}
			if address != nil {
				er.Addresses = append(er.Addresses, *address)
			}
		case types.EventTypeCoinSpent:
			address, err := parseCoinSpent(events[i].Data, events[i].Height)
			if err != nil {
				return errors.Wrap(err, "parse coin spent")
			}
			if address != nil {
				er.Addresses = append(er.Addresses, *address)
			}
		}
	}

	return nil
}

func getDecimalFromMap(m map[string]any, key string) decimal.Decimal {
	val, ok := m[key]
	if !ok {
		return decimal.Zero
	}
	str, ok := val.(string)
	if !ok {
		return decimal.Zero
	}
	str = strings.TrimSuffix(str, "utia")
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return decimal.Zero
	}
	return dec
}

func getStringFromMap(m map[string]any, key string) string {
	val, ok := m[key]
	if !ok {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func getBalanceFromMap(m map[string]any, key string) (*cosmosTypes.Coin, error) {
	val, ok := m[key]
	if !ok {
		return nil, nil
	}
	str, ok := val.(string)
	if !ok {
		return nil, nil
	}

	coin, err := cosmosTypes.ParseCoinNormalized(str)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}
