package storage

import (
	"context"

	"github.com/lib/pq"
)

var Models = []any{
	&State{},
	&Address{},
	&Block{},
	&Tx{},
	&Message{},
	&Event{},
	&Namespace{},
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Notificator interface {
	Notify(ctx context.Context, channel string, payload string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Listener interface {
	Subscribe(ctx context.Context, channels ...string) error
	Listen() chan *pq.Notification
}

const (
	ChannelHead = "head"
)
