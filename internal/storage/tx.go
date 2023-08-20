package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// ITx -
type ITx interface {
	storage.Table[*Tx]
}

// Tx -
type Tx struct {
	bun.BaseModel `bun:"tx" comment:"Table with celestia transactions." partition:"RANGE(time)"`

	Id            uint64          `bun:"id,type:bigint,pk,notnull" comment:"Unique internal id"`
	Height        uint64          `bun:",notnull" comment:"The number (height) of this block"`
	Time          time.Time       `comment:"The time of block"`
	Position      uint64          `comment:"Position in block"`
	GasWanted     uint64          `comment:"Gas wanted"`
	GasUsed       uint64          `comment:"Gas used"`
	TimeoutHeight uint64          `comment:"Block height until which the transaction is valid"`
	EventsCount   uint64          `comment:"Events count in transaction"`
	MessagesCount uint64          `comment:"Messages count in transaction"`
	Fee           decimal.Decimal `bun:",type:numeric" comment:"Paid fee"`
	Status        Status          `bun:",type:status" comment:"Transaction status"`
	Error         string          `comment:"Error string if failed"`
	Codespace     string          `comment:"Codespace"`
	Hash          []byte          `comment:"Transaction hash"`
	Memo          string          `comment:"Note or comment to send with the transaction"`

	Messages []Message `bun:"rel:has-many"`
	Events   []Event   `bun:"rel:has-many"`
}

// TableName -
func (Tx) TableName() string {
	return "tx"
}
