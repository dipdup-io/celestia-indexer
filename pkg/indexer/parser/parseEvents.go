package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func parseEvents(b types.BlockData, events []nodeTypes.Event) []storage.Event {
	result := make([]storage.Event, len(events))

	return result
}
