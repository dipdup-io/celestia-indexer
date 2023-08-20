package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IState -
type IState interface {
	storage.Table[*State]

	ByName(ctx context.Context, name string) (State, error)
}

// State -
type State struct {
	bun.BaseModel `bun:"state" comment:"Current indexer state"`

	ID         uint64    `bun:",pk,autoincrement" comment:"Unique internal identity"`
	Name       string    `bun:",unique:state_name"`
	LastHeight uint64    `comment:"Last block height"`
	LastTime   time.Time `comment:"Time of last block"`

	TotalTx            uint64 `comment:"Transactions count in celestia"`
	TotalAccounts      uint64 `comment:"Accounts count in celestia"`
	TotalNamespaces    uint64 `comment:"Namespaces count in celestia"`
	TotalNamespaceSize uint64 `comment:"Total namespace size"`
}

// TableName -
func (State) TableName() string {
	return "state"
}
