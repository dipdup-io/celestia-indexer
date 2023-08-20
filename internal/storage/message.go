package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IMessage -
type IMessage interface {
	storage.Table[*Message]
}

// Message -
type Message struct {
	bun.BaseModel `bun:"message" comment:"Table with celestia messages." partition:"RANGE(time)"`

	Id       uint64         `bun:"id,type:bigint,pk,notnull" comment:"Unique internal id"`
	Height   uint64         `bun:",notnull" comment:"The number (height) of this block"`
	Time     time.Time      `comment:"The time of block"`
	Position uint64         `comment:"Position in transaction"`
	Type     MsgType        `bun:",type:msg_type" comment:"Message type"`
	TxId     *uint64        `comment:"Parent transaction id"`
	Data     map[string]any `bun:"type:jsonb" comment:"Message data"`
}

// TableName -
func (Message) TableName() string {
	return "message"
}
