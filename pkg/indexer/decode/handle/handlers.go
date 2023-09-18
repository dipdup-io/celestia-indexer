package handle

import (
	"github.com/celestiaorg/celestia-app/pkg/namespace"
	types8 "github.com/celestiaorg/celestia-app/x/blob/types"
	types9 "github.com/celestiaorg/celestia-app/x/qgb/types"
	types7 "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	types6 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	types5 "github.com/cosmos/cosmos-sdk/x/slashing/types"
	types4 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	types2 "github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func MsgEditValidator(level types2.Level, status types.Status, m *types4.MsgEditValidator) (types.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := types.MsgEditValidator
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	if status == types.StatusFailed {
		return msgType, addresses, nil, nil
	}

	validator := storage.Validator{
		Address:           m.ValidatorAddress,
		Moniker:           m.Description.Moniker,
		Identity:          m.Description.Identity,
		Website:           m.Description.Website,
		Details:           m.Description.Details,
		Contacts:          m.Description.SecurityContact,
		Height:            uint64(level),
		Rate:              decimal.Zero,
		MinSelfDelegation: decimal.Zero,
	}

	if m.CommissionRate != nil && !m.CommissionRate.IsNil() {
		validator.Rate = decimal.RequireFromString(m.CommissionRate.String())
	}
	if m.MinSelfDelegation != nil && !m.MinSelfDelegation.IsNil() {
		validator.MinSelfDelegation = decimal.RequireFromString(m.MinSelfDelegation.String())
	}
	return msgType, addresses, &validator, err
}

func MsgBeginRedelegate(level types2.Level, m *types4.MsgBeginRedelegate) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgBeginRedelegate
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: types.MsgAddressTypeValidatorSrcAddress, address: m.ValidatorSrcAddress},
		{t: types.MsgAddressTypeValidatorDstAddress, address: m.ValidatorDstAddress},
	}, level)
	return msgType, addresses, err
}

func MsgCreateValidator(level types2.Level, status types.Status, m *types4.MsgCreateValidator) (types.MsgType, []storage.AddressWithType, *storage.Validator, error) {
	msgType := types.MsgCreateValidator
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: types.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	if status == types.StatusFailed {
		return msgType, addresses, nil, nil
	}

	validator := storage.Validator{
		Delegator:         m.DelegatorAddress,
		Address:           m.ValidatorAddress,
		Moniker:           m.Description.Moniker,
		Identity:          m.Description.Identity,
		Website:           m.Description.Website,
		Details:           m.Description.Details,
		Contacts:          m.Description.SecurityContact,
		Height:            uint64(level),
		Rate:              decimal.Zero,
		MaxRate:           decimal.Zero,
		MaxChangeRate:     decimal.Zero,
		MinSelfDelegation: decimal.Zero,
	}

	if !m.Commission.Rate.IsNil() {
		validator.Rate = decimal.RequireFromString(m.Commission.Rate.String())
	}

	if !m.Commission.MaxRate.IsNil() {
		validator.MaxRate = decimal.RequireFromString(m.Commission.MaxRate.String())
	}

	if !m.Commission.MaxChangeRate.IsNil() {
		validator.MaxChangeRate = decimal.RequireFromString(m.Commission.MaxChangeRate.String())
	}

	if !m.MinSelfDelegation.IsNil() {
		validator.MinSelfDelegation = decimal.RequireFromString(m.MinSelfDelegation.String())
	}

	return msgType, addresses, &validator, err
}

func MsgDelegate(level types2.Level, m *types4.MsgDelegate) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgDelegate
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: types.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func MsgUndelegate(level types2.Level, m *types4.MsgUndelegate) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgUndelegate
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeDelegatorAddress, address: m.DelegatorAddress},
		{t: types.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
	}, level)
	return msgType, addresses, err
}

func MsgUnjail(level types2.Level, m *types5.MsgUnjail) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgUnjail
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeValidatorAddress, address: m.ValidatorAddr},
	}, level)
	return msgType, addresses, err
}

func MsgSend(level types2.Level, m *types6.MsgSend) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgSend
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: types.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

func MsgCreateVestingAccount(level types2.Level, m *types7.MsgCreateVestingAccount) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgCreateVestingAccount
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: types.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

func MsgCreatePeriodicVestingAccount(level types2.Level, m *types7.MsgCreatePeriodicVestingAccount) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgCreatePeriodicVestingAccount
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeFromAddress, address: m.FromAddress},
		{t: types.MsgAddressTypeToAddress, address: m.ToAddress},
	}, level)
	return msgType, addresses, err
}

func MsgPayForBlobs(level types2.Level, m *types8.MsgPayForBlobs) (types.MsgType, []storage.AddressWithType, []storage.Namespace, uint64, error) {
	var blobsSize uint64
	namespaces := make([]storage.Namespace, len(m.Namespaces))

	for nsI, ns := range m.Namespaces {
		if len(m.BlobSizes) < nsI {
			return types.MsgUnknown, nil, nil, 0, errors.Errorf(
				"blob sizes length=%d is less then namespaces index=%d", len(m.BlobSizes), nsI)
		}

		appNS := namespace.Namespace{Version: ns[0], ID: ns[1:]}
		size := uint64(m.BlobSizes[nsI])
		blobsSize += size
		namespaces[nsI] = storage.Namespace{
			FirstHeight: level,
			Version:     appNS.Version,
			NamespaceID: appNS.ID,
			Size:        size,
			PfbCount:    1,
			Reserved:    appNS.IsReserved(),
		}
	}

	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeSigner, address: m.Signer},
	}, level)

	return types.MsgPayForBlobs, addresses, namespaces, blobsSize, err
}

func MsgGrantAllowance(level types2.Level, m *feegrant.MsgGrantAllowance) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgGrantAllowance
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeGranter, address: m.Granter},
		{t: types.MsgAddressTypeGrantee, address: m.Grantee},
	}, level)
	return msgType, addresses, err
}

func MsgRegisterEVMAddress(level types2.Level, m *types9.MsgRegisterEVMAddress) (types.MsgType, []storage.AddressWithType, error) {
	msgType := types.MsgRegisterEVMAddress
	addresses, err := createAddresses(addressesData{
		{t: types.MsgAddressTypeValidatorAddress, address: m.ValidatorAddress},
		// TODO: think about EVM addresses
	}, level)
	return msgType, addresses, err
}
