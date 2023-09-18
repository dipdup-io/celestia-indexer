package handle

import (
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func MsgVote(level types.Level, m *group.MsgVote) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVote
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: m.Voter},
	}, level)
	return msgType, addresses, err
}
