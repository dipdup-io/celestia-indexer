package handle

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	types2 "github.com/dipdup-io/celestia-indexer/pkg/types"
)

type addressesData []struct {
	t       types.MsgAddressType
	address string
}

func createAddresses(data addressesData, level types2.Level) ([]storage.AddressWithType, error) {
	addresses := make([]storage.AddressWithType, len(data))
	for i, d := range data {
		_, hash, err := types2.Address(d.address).Decode()
		if err != nil {
			return nil, err
		}
		addresses[i] = storage.AddressWithType{
			Type: d.t,
			Address: storage.Address{
				Hash:       hash,
				Height:     level,
				LastHeight: level,
				Address:    d.address,
				Balance:    storage.EmptyBalance(),
			},
		}
	}
	return addresses, nil
}
