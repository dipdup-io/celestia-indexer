package parser

import (
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/bytes"
	"testing"
	"time"
)

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block := types.BlockData{
		ResultBlockResults: nodeTypes.ResultBlockResults{
			TxsResults: make([]*nodeTypes.ResponseDeliverTx, 0),
		},
	}

	resultTxs := parseTxs(block)

	assert.Empty(t, resultTxs)
}

func TestParseTxs_SuccessResult(t *testing.T) {
	txRes := nodeTypes.ResponseDeliverTx{
		Code:      0,
		Data:      bytes.HexBytes{},
		Log:       "[]",
		Info:      "info",
		GasWanted: 12000,
		GasUsed:   1000,
		Events:    nil,
		Codespace: "celestia-explorer",
	}
	headerBlock := nodeTypes.Block{
		Header: nodeTypes.Header{
			Time: time.Now(),
		},
	}
	block := types.BlockData{
		ResultBlock: nodeTypes.ResultBlock{
			Block: &headerBlock,
		},
		ResultBlockResults: nodeTypes.ResultBlockResults{
			TxsResults: []*nodeTypes.ResponseDeliverTx{
				&txRes, &txRes, &txRes,
			},
		},
	}

	resultTxs := parseTxs(block)

	assert.Len(t, resultTxs, 3)
}
