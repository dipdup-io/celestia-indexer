package handle_test

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/internal/test_suite"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode"
	"github.com/fatih/structs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createMsgVote() types.Msg {
	// Data from: 0A4BA0A30449C3269F313B5D974560F8D3A8179BE994054724898FF2D6866928
	m := group.MsgVote{
		ProposalId: 1,
		Voter:      "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
		Option:     group.VOTE_OPTION_YES,
		Metadata:   "",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgVote(t *testing.T) {
	m := createMsgVote()
	blob, now := testsuite.EmptyBlock()
	position := 7

	dm, err := decode.Message(m, blob.Height, blob.Block.Time, position, storageTypes.StatusSuccess)

	addressesExpected := []storage.AddressWithType{
		{
			Type: storageTypes.MsgAddressTypeVoter,
			Address: storage.Address{
				Id:         0,
				Height:     blob.Height,
				LastHeight: blob.Height,
				Address:    "celestia1prxtghtsjrdwdtkt82kye3a7yukmcay6x9uyts",
				Hash:       []byte{8, 204, 180, 93, 112, 144, 218, 230, 174, 203, 58, 172, 76, 199, 190, 39, 45, 188, 116, 154},
				Balance: storage.Balance{
					Total: decimal.Zero,
				},
			},
		},
	}

	msgExpected := storage.Message{
		Id:        0,
		Height:    blob.Height,
		Time:      now,
		Position:  7,
		Type:      storageTypes.MsgVote,
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
