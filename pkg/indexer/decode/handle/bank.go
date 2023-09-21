package handle

import (
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgSend represents a message to send coins from one account to another.
func MsgSend(level types.Level, m *cosmosBankTypes.MsgSend) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSend
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

// MsgMultiSend represents an arbitrary multi-in, multi-out send message.
func MsgMultiSend(level types.Level, m *cosmosBankTypes.MsgMultiSend) (storageTypes.MsgType, []storage.AddressWithType, error) {

}
