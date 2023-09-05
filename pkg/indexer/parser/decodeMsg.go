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
	"github.com/shopspring/decimal"
)

type decodedMsg struct {
	msg       storage.Message
	blobsSize uint64
	addresses []storage.AddressWithType
}

func decodeMsg(b types.BlockData, msg cosmosTypes.Msg, position int) (d decodedMsg, err error) {
	d.msg.Height = b.Height
	d.msg.Time = b.Block.Time
	d.msg.Position = uint64(position)
	d.msg.Data = structs.Map(msg)

	switch msg.(type) {
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		d.msg.Type, d.addresses, err = handledMsgWithdrawValidatorCommission(b, msg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.msg.Type = storageTypes.MsgTypeWithdrawDelegatorReward
	case *cosmosStakingTypes.MsgEditValidator:
		d.msg.Type = storageTypes.MsgTypeEditValidator
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.msg.Type = storageTypes.MsgTypeBeginRedelegate
	case *cosmosStakingTypes.MsgCreateValidator:
		d.msg.Type = storageTypes.MsgTypeCreateValidator
	case *cosmosStakingTypes.MsgDelegate:
		d.msg.Type = storageTypes.MsgTypeDelegate
	case *cosmosStakingTypes.MsgUndelegate:
		d.msg.Type = storageTypes.MsgTypeUndelegate
	case *cosmosSlashingTypes.MsgUnjail:
		d.msg.Type = storageTypes.MsgTypeUnjail
	case *cosmosBankTypes.MsgSend:
		d.msg.Type = storageTypes.MsgTypeSend
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.msg.Type = storageTypes.MsgTypeCreateVestingAccount
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.msg.Type = storageTypes.MsgTypeCreatePeriodicVestingAccount
	case *appBlobTypes.MsgPayForBlobs:
		d.msg.Type = storageTypes.MsgTypePayForBlobs
		d.msg.Namespace, d.blobsSize, err = handlePfb(b, msg)
	case *cosmosFeegrant.MsgGrantAllowance:
		d.msg.Type = storageTypes.MsgTypeGrantAllowance
	default:
		d.msg.Type = storageTypes.MsgTypeUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	return
}

func handledMsgWithdrawValidatorCommission(b types.BlockData, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeWithdrawValidatorCommission

	m, ok := msg.(*cosmosDistributionTypes.MsgWithdrawValidatorCommission)
	if !ok {
		return storageTypes.MsgTypeUnknown, nil, errors.Errorf("error on decoding '%T' in cosmosDistributionTypes.MsgWithdrawValidatorCommission", msg)
	}

	addresses := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Height:  b.Height,
				Hash:    []byte(m.ValidatorAddress),
				Balance: decimal.Zero,
			},
		},
	}

	return msgType, addresses, nil
}

func handlePfb(b types.BlockData, msg cosmosTypes.Msg) ([]storage.Namespace, uint64, error) {
	payForBlobsMsg, ok := msg.(*appBlobTypes.MsgPayForBlobs)
	if !ok {
		return nil, 0, errors.Errorf("error on decoding '%T' in appBlobTypes.MsgPayForBlobs", msg)
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
