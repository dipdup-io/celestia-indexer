package postgres

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

const (
	connectionName       = "celestia_notifications"
	minReconnectInterval = 10 * time.Second
	maxReconnectInterval = 20 * time.Second
)

type Notificator struct {
	db *bun.DB
	l  *pq.Listener
}

func NewNotificator(db *bun.DB) *Notificator {
	return &Notificator{
		l: pq.NewListener(
			connectionName,
			minReconnectInterval,
			maxReconnectInterval,
			nil,
		),
		db: db,
	}
}

func (n *Notificator) Notify(ctx context.Context, channel string, payload string) error {
	_, err := n.db.ExecContext(ctx, "NOTIFY ?, ?", bun.Ident(channel), payload)
	return err
}

func (n *Notificator) Listen() chan *pq.Notification {
	return n.l.Notify
}

func (n *Notificator) Subscribe(ctx context.Context, channels ...string) error {
	for i := range channels {
		if err := n.l.Listen(channels[i]); err != nil {
			return err
		}
	}
	return nil
}

func (n *Notificator) Close() error {
	return n.l.Close()
}
