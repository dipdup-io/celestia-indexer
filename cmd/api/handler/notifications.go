package handler

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/websocket"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

func ListenNotifications(ctx context.Context, manager *websocket.Manager, listener storage.Listener) {
	if err := listener.Subscribe(ctx, "head"); err != nil {
		log.Err(err).Msg("subscribe on notifications")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-listener.Listen():
			if !ok {
				continue
			}
			switch msg.Channel {
			case "head":
				var block storage.Block
				if err := json.Unmarshal([]byte(msg.Extra), &block); err != nil {
					log.Err(err).Msg("block decoding in notifications")
					continue
				}
				manager.NotifyAll(block)
			default:
				log.Warn().Str("channel", msg.Channel).Msg("unknown channel")
			}
		}

	}
}
