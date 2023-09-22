package handle

import (
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

// MsgSubmitProposal defines a sdk.Msg type that supports submitting arbitrary
// proposal Content.
func MsgSubmitProposal(level types.Level, proposerAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitProposal
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeProposer, address: proposerAddress},
	}, level)
	return msgType, addresses, err
}

// MsgExecLegacyContent is used to wrap the legacy content field into a message.
// This ensures backwards compatibility with v1beta1.MsgSubmitProposal.
func MsgExecLegacyContent(level types.Level, m *cosmosGovTypesV1.MsgExecLegacyContent) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgExecLegacyContent
	addresses, err := createAddresses(
		addressesData{
			{t: storageTypes.MsgAddressTypeAuthority, address: m.Authority},
		}, level)
	return msgType, addresses, err
}

func MsgVote(level types.Level, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVote
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, level)
	return msgType, addresses, err
}

func MsgVoteWeighted(level types.Level, voterAddress string) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgVoteWeighted
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeVoter, address: voterAddress},
	}, level)
	return msgType, addresses, err
}
