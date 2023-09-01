package parser

import (
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/bytes"
	"testing"
)

func TestParseEvents_EmptyEventsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: nodeTypes.ResultBlockResults{
			TxsResults: make([]*nodeTypes.ResponseDeliverTx, 0),
		},
	}

	resultEvents := parseEvents(block, make([]nodeTypes.Event, 0), nil)

	assert.Empty(t, resultEvents)
}

func TestParseEvents_SuccessTx(t *testing.T) {
	events := []nodeTypes.Event{
		{
			Type: "coin_spent",
			Attributes: []nodeTypes.EventAttribute{
				{
					Key:   bytes.HexBytes("c3BlbmRlcg=="),
					Value: bytes.HexBytes("Y2VsZXN0aWExdjY5bnB6NncwN3h0NGhkdWU5eGR3a3V4eHZ2ZDZlYTl5MjZlcXI="),
					Index: true,
				},
				{
					Key:   bytes.HexBytes("YW1vdW50"),
					Value: bytes.HexBytes("NzAwMDB1dGlh"),
					Index: true,
				},
			},
		},
	}

	txRes := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      bytes.HexBytes{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    events,
		Codespace: "celestia-explorer",
	}
	block, now := createBlock(txRes, 1)

	var txId *uint64
	resultEvents := parseEvents(block, events, txId)

	assert.Len(t, resultEvents, 1)

	e := resultEvents[0]
	assert.Equal(t, block.Height, e.Height)
	assert.Equal(t, now, e.Time)
	assert.Equal(t, uint64(0), e.Position)
	assert.Equal(t, storageTypes.EventTypeCoinSpent, e.Type)
	assert.Equal(t, txId, e.TxId)

	attrs := map[string]any{
		bytes.HexBytes("c3BlbmRlcg==").String(): bytes.HexBytes("Y2VsZXN0aWExdjY5bnB6NncwN3h0NGhkdWU5eGR3a3V4eHZ2ZDZlYTl5MjZlcXI="),
		bytes.HexBytes("YW1vdW50").String():     bytes.HexBytes("NzAwMDB1dGlh"),
	}
	assert.Equal(t, attrs, e.Data)
}
