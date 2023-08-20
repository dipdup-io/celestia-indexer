package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IBlock -
type IBlock interface {
	storage.Table[*Block]

	Last(ctx context.Context) (Block, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"table:block" comment:"Table with celestia blocks."`

	Height       uint64    `bun:",pk,notnull" comment:"The number (height) of this block"`
	Time         time.Time `comment:"The time of block"`
	VersionBlock string    `comment:"Block version"`
	VersionApp   string    `comment:"App version"`

	TxCount uint64 `comment:"Count of transactions in block"`

	Hash               []byte `comment:"Block hash"`
	ParentHash         []byte `comment:"Hash of parent block"`
	LastCommitHash     []byte `comment:"Last commit hash"`
	DataHash           []byte `comment:"Data hash"`
	ValidatorsHash     []byte `comment:"Validators hash"`
	NextValidatorsHash []byte `comment:"Next validators hash"`
	ConsensusHash      []byte `comment:"Consensus hash"`
	AppHash            []byte `comment:"App hash"`
	LastResultsHash    []byte `comment:"Last results hash"`
	EvidenceHash       []byte `comment:"Evidence hash"`
	ProposerAddress    []byte `comment:"Proposer address"`

	Txs []Tx `bun:"rel:has-many"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
