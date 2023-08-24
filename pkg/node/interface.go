package node

import (
	"context"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
)

type API interface {
	GetHead(ctx context.Context) (types.ResultBlock, error)
	GetBlock(ctx context.Context, level storage.Level) (types.ResultBlock, error)
	GetGenesis(ctx context.Context) (types.Genesis, error)
}
