package decode

import (
	"testing"

	"cosmossdk.io/math"
	cosmosCodecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgEditValidator(t *testing.T) {
	m := createMsgEditValidator()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgEditValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
		Validator: &storage.Validator{
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			Moniker:           "newAgeValidator",
			Identity:          "UPort:1",
			Website:           "https://google.com",
			Contacts:          "tryme@gmail.com",
			Details:           "trust",
			Rate:              decimal.Zero,
			MinSelfDelegation: decimal.Zero,
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgBeginRedelegate

func createMsgBeginRedelegate() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorSrcAddress: "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		ValidatorDstAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgBeginRedelegate(t *testing.T) {
	m := createMsgBeginRedelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:       []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorSrcAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorDstAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:       []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgBeginRedelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgCreateValidator

func createMsgCreateValidator() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgCreateValidator{
		Description:       cosmosStakingTypes.Description{},
		Commission:        cosmosStakingTypes.CommissionRates{},
		MinSelfDelegation: cosmosTypes.NewInt(1),
		DelegatorAddress:  "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
		ValidatorAddress:  "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
		Pubkey:            nil,
		Value:             cosmosTypes.Coin{},
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgCreateValidator(t *testing.T) {
	m := createMsgCreateValidator()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
				Hash:       []byte{0x74, 0x2b, 0x74, 0xc3, 0xe7, 0xbf, 0xc9, 0xf5, 0xc4, 0xe1, 0x5d, 0xa9, 0x89, 0x97, 0x83, 0xea, 0x9f, 0xf, 0xf1, 0x49},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
				Hash:       []byte{0x4a, 0xb, 0xf8, 0x99, 0x89, 0xe7, 0xa, 0xe3, 0xf, 0x3b, 0x62, 0x94, 0xa2, 0x97, 0x5a, 0x17, 0x5c, 0x48, 0x22, 0xa3},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgCreateValidator,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
		Validator: &storage.Validator{
			Delegator:         "celestia1ws4hfsl8hlylt38ptk5cn9ura20slu2fnkre76",
			Address:           "celestiavaloper1fg9l3xvfuu9wxremv2229966zawysg4r40gw5x",
			Rate:              decimal.Zero,
			MaxRate:           decimal.Zero,
			MaxChangeRate:     decimal.Zero,
			MinSelfDelegation: decimal.RequireFromString("1"),
			Height:            uint64(blob.Height),
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgDelegate

func createMsgDelegate() cosmosTypes.Msg {

	msgDelegate := cosmosStakingTypes.MsgDelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
		Amount: cosmosTypes.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1000),
		},
	}

	return &msgDelegate
}

func TestDecodeMsg_SuccessOnMsgDelegate(t *testing.T) {
	msgDelegate := createMsgDelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(msgDelegate, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper12c6cwd0kqlg48sdhjnn9f0z82g0c82fmrl7j9y",
				Hash:       []byte{0x56, 0x35, 0x87, 0x35, 0xf6, 0x7, 0xd1, 0x53, 0xc1, 0xb7, 0x94, 0xe6, 0x54, 0xbc, 0x47, 0x52, 0x1f, 0x83, 0xa9, 0x3b},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgDelegate,
		TxId:      0,
		Data:      structs.Map(msgDelegate),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUndelegate

func createMsgUndelegate() cosmosTypes.Msg {
	m := cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
		ValidatorAddress: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
		Amount: cosmosTypes.Coin{
			Denom:  "utia",
			Amount: math.NewInt(1001),
		},
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUndelegate(t *testing.T) {
	m := createMsgUndelegate()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeDelegatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vysgwc9mykfz5249g9thjlffx6nha0kkwsvs37",
				Hash:       []byte{0x61, 0x20, 0x87, 0x60, 0xbb, 0x25, 0x92, 0x2a, 0x2a, 0xa5, 0x41, 0x57, 0x79, 0x7d, 0x29, 0x36, 0xa7, 0x7e, 0xbe, 0xd6},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUndelegate,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgUnjail

func createMsgUnjail() cosmosTypes.Msg {
	m := cosmosSlashingTypes.MsgUnjail{
		ValidatorAddr: "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
	}
	return &m
}

func TestDecodeMsg_SuccessOnMsgUnjail(t *testing.T) {
	m := createMsgUnjail()
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeValidatorAddress,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestiavaloper170qq26qenw420ufd5py0r59kpg3tj2m7dqkpym",
				Hash:       []byte{0xf3, 0xc0, 0x5, 0x68, 0x19, 0x9b, 0xaa, 0xa7, 0xf1, 0x2d, 0xa0, 0x48, 0xf1, 0xd0, 0xb6, 0xa, 0x22, 0xb9, 0x2b, 0x7e},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnjail,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
}

// MsgGrantAllowance

func createMsgGrantAllowance() cosmosTypes.Msg {
	m := cosmosFeegrant.MsgGrantAllowance{
		Granter:   "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
		Grantee:   "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
		Allowance: cosmosCodecTypes.UnsafePackAny(cosmosFeegrant.BasicAllowance{}),
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgGrantAllowance(t *testing.T) {
	m := createMsgGrantAllowance()
	blob, now := testsuite.EmptyBlock()
	position := 4

	dm, err := Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeGranter,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia18r6ujzzkg6ku9sr39nxy4847q4qea5kg4a8pxv",
				Hash:       []byte{0x38, 0xf5, 0xc9, 0x8, 0x56, 0x46, 0xad, 0xc2, 0xc0, 0x71, 0x2c, 0xcc, 0x4a, 0x9e, 0xbe, 0x5, 0x41, 0x9e, 0xd2, 0xc8},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
		{
			Type: storageTypes.MsgAddressTypeGrantee,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1vnflc6322f8z7cpl28r7un5dxhmjxghc20aydq",
				Hash:       []byte{0x64, 0xd3, 0xfc, 0x6a, 0x2a, 0x52, 0x4e, 0x2f, 0x60, 0x3f, 0x51, 0xc7, 0xee, 0x4e, 0x8d, 0x35, 0xf7, 0x23, 0x22, 0xf8},
				Balance: storage.Balance{
					Id:    0,
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  4,
		Type:      storageTypes.MsgGrantAllowance,
		TxId:      0,
		Data:      structs.Map(m),
		Namespace: nil,
		Addresses: addressesExpected,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
	assert.Equal(t, addressesExpected, dm.Addresses)
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
	blob, now := testsuite.EmptyBlock()
	position := 0

	dm, err := Message(msgUnknown, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnknown,
		TxId:      0,
		Data:      structs.Map(msgUnknown),
		Namespace: nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(0), dm.BlobsSize)
	assert.Equal(t, msgExpected, dm.Msg)
}
