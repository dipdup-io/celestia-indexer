package parser

import (
	"cosmossdk.io/math"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	cosmosCodecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MsgWithdrawValidatorCommission

func createMsgWithdrawValidatorCommission() cosmosTypes.Msg {
	m := cosmosDistributionTypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgWithdrawValidatorCommission(t *testing.T) {
	m := createMsgWithdrawValidatorCommission()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeWithdrawValidatorCommission,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgWithdrawDelegatorReward

func createMsgWithdrawDelegatorReward() cosmosTypes.Msg {
	m := cosmosDistributionTypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgWithdrawDelegatorReward(t *testing.T) {
	m := createMsgWithdrawDelegatorReward()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeWithdrawDelegatorReward,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgEditValidator

func createMsgEditValidator() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgEditValidator{
		Description: cosmosStakingTypes.Description{
			Moniker:         "newAgeValidator",
			Identity:        "UPort:1",
			Website:         "https://google.com",
			SecurityContact: "tryme@gmail.com",
			Details:         "trust",
		},
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		CommissionRate:    nil,
		MinSelfDelegation: nil,
		EvmAddress:        "0x10E0271ec47d55511a047516f2a7301801d55eaB",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgEditValidator(t *testing.T) {
	m := createMsgEditValidator()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeEditValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgBeginRedelegate

func createMsgBeginRedelegate() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorSrcAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		ValidatorDstAddress: "celestiavaloper1fg9l3xvfuu9wxremv2288777zawysg4r40gw7x",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgBeginRedelegate(t *testing.T) {
	m := createMsgBeginRedelegate()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeBeginRedelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorSrcAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorDstAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2288777zawysg4r40gw7x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgCreateValidator

func createMsgCreateValidator() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgCreateValidator{
		Description:       cosmosStakingTypes.Description{},
		Commission:        cosmosStakingTypes.CommissionRates{},
		MinSelfDelegation: cosmosTypes.Int{}, // nolint
		DelegatorAddress:  "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre77",
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw7x",
		Pubkey:            nil,
		Value:             cosmosTypes.Coin{},
		EvmAddress:        "",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateValidator(t *testing.T) {
	m := createMsgCreateValidator()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeCreateValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre77"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw7x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgDelegate

func createMsgDelegate() cosmosTypes.Msg {

	msgDelegate := cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: "celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp95kth",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g77x",
		Amount: cosmosTypes.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1000),
		},
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgDelegate(t *testing.T) {
	msgDelegate := createMsgDelegate()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, msgDelegate, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeDelegate,
		TxId:      0,
		Data:      structs.Map(msgDelegate),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp95kth"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g77x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgUndelegate

func createMsgUndelegate() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp99kth",
		ValidatorAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g88x",
		Amount: cosmosTypes.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1001),
		},
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUndelegate(t *testing.T) {
	m := createMsgUndelegate()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeUndelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1A2kqw44hdq5dwlcvsw8f2l49lkehtf9wp99kth"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g88x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgUnjail

func createMsgUnjail() cosmosTypes.Msg {
	m := cosmosSlashingTypes.MsgUnjail{
		ValidatorAddr: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g11x",
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUnjail(t *testing.T) {
	m := createMsgUnjail()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeUnjail,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40g11x"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgSend

func createMsgSend() cosmosTypes.Msg {
	m := cosmosBankTypes.MsgSend{
		FromAddress: "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: cosmosTypes.Coins{
			cosmosTypes.Coin{
				Denom:  "utia",
				Amount: math.NewInt(1000),
			},
		},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSend(t *testing.T) {
	msgSend := createMsgSend()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, msgSend, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeSend,
		TxId:      0,
		Data:      structs.Map(msgSend),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeFromAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgCreateVestingAccount

func createMsgCreateVestingAccount() cosmosTypes.Msg {
	m := cosmosVestingTypes.MsgCreateVestingAccount{
		FromAddress: "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
		ToAddress:   "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		Amount: cosmosTypes.Coins{
			cosmosTypes.Coin{
				Denom:  "utia",
				Amount: math.NewInt(1000),
			},
		},
		EndTime: 0,
		Delayed: false,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateVestingAccount(t *testing.T) {
	m := createMsgCreateVestingAccount()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeCreateVestingAccount,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeFromAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgCreatePeriodicVestingAccount

func createMsgCreatePeriodicVestingAccount() cosmosTypes.Msg {
	m := cosmosVestingTypes.MsgCreatePeriodicVestingAccount{
		FromAddress:    "celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7",
		ToAddress:      "celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l",
		StartTime:      0,
		VestingPeriods: nil,
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreatePeriodicVestingAccount(t *testing.T) {
	msgCreatePeriodicVestingAccount := createMsgCreatePeriodicVestingAccount()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, msgCreatePeriodicVestingAccount, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeCreatePeriodicVestingAccount,
		TxId:      0,
		Data:      structs.Map(msgCreatePeriodicVestingAccount),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeFromAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia10a0qvvg53svyfvmuf5azx779xrpwn9lxzlfkn7"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeToAddress,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vsvx8n7f8dh5udesqqhgrjutyun7zqrgehdq2l"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgPayForBlob

func createMsgPayForBlob() cosmosTypes.Msg {

	msgPayForBlob := appBlobTypes.MsgPayForBlobs{
		Signer:           "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		Namespaces:       [][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22}},
		BlobSizes:        []uint32{1},
		ShareCommitments: [][]byte{{176, 28, 134, 119, 32, 117, 87, 107, 231, 67, 121, 255, 209, 106, 52, 99, 88, 183, 85, 36, 67, 137, 98, 199, 144, 159, 13, 178, 111, 190, 121, 36}},
		ShareVersions:    []uint32{0},
	}

	return &msgPayForBlob
}

func TestDecodeMsg_SuccessOnPayForBlob(t *testing.T) {
	msgPayForBlob := createMsgPayForBlob()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, msgPayForBlob, position)

	msgExpected := storage.Message{
		Id:       0,
		Height:   blob.Height,
		Time:     now,
		Position: 0,
		Type:     storageTypes.MsgTypePayForBlobs,
		TxId:     0,
		Data:     structs.Map(msgPayForBlob),
		Namespace: []storage.Namespace{
			{
				Id:          0,
				FirstHeight: blob.Height,
				Version:     0,
				NamespaceID: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 189, 44, 204, 197, 144, 206, 197, 121, 37, 22},
				Size:        1,
				PfbCount:    1,
				Reserved:    false,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
}

// MsgGrantAllowance

func createMsgGrantAllowance() cosmosTypes.Msg {
	m := cosmosFeegrant.MsgGrantAllowance{
		Granter:   "celestia1l9qjhhnxc0t6tt93q8396gu0vttwlcc238gyvr",
		Grantee:   "celestia1vut644llcgwyvysmma6ww2xkmdytc8xspty5kx",
		Allowance: cosmosCodecTypes.UnsafePackAny(cosmosFeegrant.BasicAllowance{}),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgGrantAllowance(t *testing.T) {
	m := createMsgGrantAllowance()
	blob, now := createEmptyBlock()
	position := 4

	dm, err := decodeMsg(blob, m, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgTypeGrantAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
	}

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.TxAddressTypeGranter,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1l9qjhhnxc0t6tt93q8396gu0vttwlcc238gyvr"),
				Balance: decimal.Zero,
			},
		},
		{
			Type: storageTypes.TxAddressTypeGrantee,
			Address: storage.Address{
				Id:      0,
				Height:  blob.Height,
				Hash:    []byte("celestia1vut644llcgwyvysmma6ww2xkmdytc8xspty5kx"),
				Balance: decimal.Zero,
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
	assert.Equal(t, addressesExpected, dm.addresses)
}

// MsgUnknown

type UnknownMsgType struct{}

func (u *UnknownMsgType) Reset()                               {}
func (u *UnknownMsgType) String() string                       { return "unknown" }
func (u *UnknownMsgType) ProtoMessage()                        {}
func (u *UnknownMsgType) ValidateBasic() error                 { return nil }
func (u *UnknownMsgType) GetSigners() []cosmosTypes.AccAddress { return nil }

func createMsgUnknown() cosmosTypes.Msg {
	msgUnknown := UnknownMsgType{}
	return &msgUnknown
}

func TestDecodeMsg_MsgUnknown(t *testing.T) {
	msgUnknown := createMsgUnknown()
	blob, now := createEmptyBlock()
	position := 0

	dm, err := decodeMsg(blob, msgUnknown, position)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgTypeUnknown,
		TxId:      0,
		Data:      structs.Map(msgUnknown),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.blobsSize)
	assert.Equal(t, msgExpected, dm.msg)
}
