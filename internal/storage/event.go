package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IEvent -
type IEvent interface {
	storage.Table[*Event]
}

// Event -
type Event struct {
	bun.BaseModel `bun:"event" comment:"Table with celestia events." partition:"RANGE(time)"`

	Id       uint64         `bun:"id,type:bigint,pk,notnull" comment:"Unique internal id"`
	Height   uint64         `bun:",notnull" comment:"The number (height) of this block"`
	Time     time.Time      `comment:"The time of block"`
	Position uint64         `comment:"Position in transaction"`
	Type     EventType      `bun:",type:event_type" comment:"Event type"`
	TxId     *uint64        `comment:"Transaction id"`
	Data     map[string]any `bun:"type:jsonb" comment:"Event data"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}
