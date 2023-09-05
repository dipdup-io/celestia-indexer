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
		d.msg.Type, d.addresses, err = handleMsgWithdrawValidatorCommission(b.Height, msg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.msg.Type, d.addresses, err = handleMsgWithdrawDelegatorReward(b.Height, msg)
	case *cosmosStakingTypes.MsgEditValidator:
		d.msg.Type, d.addresses, err = handleMsgEditValidator(b.Height, msg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.msg.Type, d.addresses, err = handleMsgBeginRedelegate(b.Height, msg)
	case *cosmosStakingTypes.MsgCreateValidator:
		d.msg.Type, d.addresses, err = handleMsgCreateValidator(b.Height, msg)
	case *cosmosStakingTypes.MsgDelegate:
		d.msg.Type, d.addresses, err = handleMsgDelegate(b.Height, msg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.msg.Type, d.addresses, err = handleMsgUndelegate(b.Height, msg)
	case *cosmosSlashingTypes.MsgUnjail:
		d.msg.Type, d.addresses, err = handleMsgUnjail(b.Height, msg)
	case *cosmosBankTypes.MsgSend:
		d.msg.Type, d.addresses, err = handleMsgSend(b.Height, msg)
	case *cosmosVestingTypes.MsgCreateVestingAccount: // last
		d.msg.Type, d.addresses, err = handleMsgCreateVestingAccount(b.Height, msg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.msg.Type, d.addresses, err = handleMsgCreatePeriodicVestingAccount(b.Height, msg)
	case *appBlobTypes.MsgPayForBlobs:
		d.msg.Namespace, d.blobsSize, err = handleMsgPayForBlobs(b, msg)
		d.msg.Type = storageTypes.MsgTypePayForBlobs
	case *cosmosFeegrant.MsgGrantAllowance:
		d.msg.Type, d.addresses, err = handleMsgGrantAllowance(b.Height, msg)
	default:
		d.msg.Type = storageTypes.MsgTypeUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	return
}

type addressesData []struct {
	t       storageTypes.TxAddressType
	address string
}

func createAddresses(data addressesData, level storage.Level) []storage.AddressWithType {
	addresses := make([]storage.AddressWithType, len(data))
	for i, d := range data {
		addresses[i] = storage.AddressWithType{
			Type: d.t,
			Address: storage.Address{
				Height:  level,
				Hash:    []byte(d.address),
				Balance: decimal.Zero,
			},
		}
	}
	return addresses
}

func handleMsgWithdrawValidatorCommission(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeWithdrawValidatorCommission
	m := msg.(*cosmosDistributionTypes.MsgWithdrawValidatorCommission)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgWithdrawDelegatorReward(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeWithdrawDelegatorReward
	m := msg.(*cosmosDistributionTypes.MsgWithdrawDelegatorReward)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)

	return msgType, addresses, nil
}

func handleMsgEditValidator(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeEditValidator
	m := msg.(*cosmosStakingTypes.MsgEditValidator)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgBeginRedelegate(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeBeginRedelegate
	m := msg.(*cosmosStakingTypes.MsgBeginRedelegate)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorSrcAddress, address: m.ValidatorSrcAddress},
		{t: storageTypes.TxAddressTypeValidatorDstAddress, address: m.ValidatorDstAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgCreateValidator(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeCreateValidator
	m := msg.(*cosmosStakingTypes.MsgCreateValidator)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgDelegate(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeDelegate
	m := msg.(*cosmosStakingTypes.MsgDelegate)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgUndelegate(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeUndelegate
	m := msg.(*cosmosStakingTypes.MsgUndelegate)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgUnjail(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeUnjail
	m := msg.(*cosmosSlashingTypes.MsgUnjail)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeValidatorAddress, address: m.ValidatorAddr},
	}, level)
	return msgType, addresses, nil
}

func handleMsgSend(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeSend
	m := msg.(*cosmosBankTypes.MsgSend)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.TxAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgCreateVestingAccount(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeCreateVestingAccount
	m := msg.(*cosmosVestingTypes.MsgCreateVestingAccount)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.TxAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgCreatePeriodicVestingAccount(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeCreatePeriodicVestingAccount
	m := msg.(*cosmosVestingTypes.MsgCreatePeriodicVestingAccount)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeFromAddress, address: m.FromAddress},
		{t: storageTypes.TxAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, nil
}

func handleMsgPayForBlobs(b types.BlockData, msg cosmosTypes.Msg) ([]storage.Namespace, uint64, error) {
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

func handleMsgGrantAllowance(level storage.Level, msg cosmosTypes.Msg) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgTypeGrantAllowance
	m := msg.(*cosmosFeegrant.MsgGrantAllowance)
	addresses := createAddresses(addressesData{
		{t: storageTypes.TxAddressTypeGranter, address: m.Granter},
		{t: storageTypes.TxAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, nil
}
