package rpc

import (
	"context"
	"fmt"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

type getBlockResult struct {
	Result types.ResultBlock `json:"result"`
}

func (api *API) GetBlock(ctx context.Context, level storage.Level) (types.ResultBlock, error) {
	var search string
	if level != 0 {
		search = fmt.Sprintf("?height=%d", level)
	}

	path := fmt.Sprintf("block%s", search)

	var gbr getBlockResult
	if err := api.get(ctx, path, &gbr); err != nil {
		api.log.Err(err).Msg("node get block request")
		return types.ResultBlock{}, err
	}

	return gbr.Result, nil
}
