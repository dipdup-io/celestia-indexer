package parser

import (
	"github.com/celestiaorg/celestia-app/pkg/namespace"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

func decodeMsg(b types.BlockData, msg cosmosTypes.Msg, position int) (sMsg storage.Message, blobsSize uint64, err error) {
	sMsg.Height = b.Height
	sMsg.Time = b.Block.Time
	sMsg.Position = uint64(position)
	sMsg.Data = structs.Map(msg)

	switch msg.(type) {
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		sMsg.Type = storageTypes.MsgTypeWithdrawValidatorCommission
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		sMsg.Type = storageTypes.MsgTypeWithdrawDelegatorReward
	case *cosmosStakingTypes.MsgEditValidator:
		sMsg.Type = storageTypes.MsgTypeEditValidator
	case *cosmosStakingTypes.MsgBeginRedelegate:
		sMsg.Type = storageTypes.MsgTypeBeginRedelegate
	case *cosmosStakingTypes.MsgCreateValidator:
		sMsg.Type = storageTypes.MsgTypeCreateValidator
	case *cosmosStakingTypes.MsgDelegate:
		sMsg.Type = storageTypes.MsgTypeDelegate
	case *cosmosStakingTypes.MsgUndelegate:
		sMsg.Type = storageTypes.MsgTypeUndelegate
	case *cosmosSlashingTypes.MsgUnjail:
		sMsg.Type = storageTypes.MsgTypeUnjail
	case *cosmosBankTypes.MsgSend:
		sMsg.Type = storageTypes.MsgTypeSend
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		sMsg.Type = storageTypes.MsgTypeCreateVestingAccount
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		sMsg.Type = storageTypes.MsgTypeCreatePeriodicVestingAccount
	case *appBlobTypes.MsgPayForBlobs:
		sMsg.Type = storageTypes.MsgTypePayForBlobs
		sMsg.Namespace, blobsSize, err = handlePfb(b, msg)
	case *cosmosFeegrant.MsgGrantAllowance:
		sMsg.Type = storageTypes.MsgTypeGrantAllowance
	default:
		sMsg.Type = storageTypes.MsgTypeUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	return
}

func handlePfb(b types.BlockData, msg cosmosTypes.Msg) ([]storage.Namespace, uint64, error) {
	payForBlobsMsg, ok := msg.(*appBlobTypes.MsgPayForBlobs)
	if !ok {
		return nil, 0, errors.Errorf("error on decoding %T", msg)
	}

	var blobsSize uint64
	namespaces := make([]storage.Namespace, len(payForBlobsMsg.Namespaces))

	for nsI, ns := range payForBlobsMsg.Namespaces {
		if len(payForBlobsMsg.BlobSizes) < nsI {
			return nil, 0, errors.Errorf(
				"blob sizes length=%d is less then namespaces index=%d", len(payForBlobsMsg.BlobSizes), nsI)
		}

		appNS := namespace.Namespace{Version: ns[0], ID: ns[1:]}
		size := uint64(payForBlobsMsg.BlobSizes[nsI])
		blobsSize += size
		namespaces[nsI] = storage.Namespace{
			FirstHeight: b.Height,
			Version:     appNS.Version,
			NamespaceID: appNS.ID,
			Size:        size,
			PfbCount:    1,
			Reserved:    appNS.IsReserved(),
		}
	}

	return namespaces, blobsSize, nil
}
