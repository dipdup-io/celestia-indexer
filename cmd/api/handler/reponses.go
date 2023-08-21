package handler

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Address struct {
	Id      uint64 `json:"id"`
	Height  uint64 `json:"first_height"`
	Balance string `json:"balance"`
	Hash    string `json:"hash"`
}

func NewAddress(addr storage.Address) (Address, error) {
	hash, err := EncodeAddress(addr.Hash)
	if err != nil {
		return Address{}, err
	}
	return Address{
		Id:      addr.Id,
		Height:  addr.Height,
		Balance: addr.Balance.String(),
		Hash:    hash,
	}, nil
}

type Block struct {
	Id                 uint64    `json:"id"`
	Height             uint64    `json:"height"`
	Time               time.Time `json:"time"`
	VersionBlock       string    `json:"version_block"`
	VersionApp         string    `json:"version_app"`
	TxCount            uint64    `json:"tx_count"`
	Hash               []byte    `json:"hash"`
	ParentHash         []byte    `json:"parent_hash"`
	LastCommitHash     []byte    `json:"last_commit_hash"`
	DataHash           []byte    `json:"data_hash"`
	ValidatorsHash     []byte    `json:"validators_hash"`
	NextValidatorsHash []byte    `json:"next_validators_hash"`
	ConsensusHash      []byte    `json:"consensus_hash"`
	AppHash            []byte    `json:"app_hash"`
	LastResultsHash    []byte    `json:"last_results_hash"`
	EvidenceHash       []byte    `json:"evidence_hash"`
	ProposerAddress    []byte    `json:"proposer_address"`
}

func NewBlock(block storage.Block) Block {
	return Block{
		Id:                 block.Id,
		Height:             block.Height,
		Time:               block.Time,
		VersionBlock:       block.VersionBlock,
		VersionApp:         block.VersionApp,
		TxCount:            block.TxCount,
		Hash:               block.Hash,
		ParentHash:         block.ParentHash,
		LastCommitHash:     block.LastCommitHash,
		DataHash:           block.DataHash,
		ValidatorsHash:     block.ValidatorsHash,
		NextValidatorsHash: block.NextValidatorsHash,
		ConsensusHash:      block.ConsensusHash,
		AppHash:            block.AppHash,
		LastResultsHash:    block.LastResultsHash,
		EvidenceHash:       block.EvidenceHash,
		ProposerAddress:    block.ProposerAddress,
	}
}

type Event struct {
	Id       uint64         `json:"id"`
	Height   uint64         `json:"height"`
	Time     time.Time      `json:"time"`
	Position uint64         `json:"position"`
	Type     string         `json:"type"`
	TxId     *uint64        `json:"tx_id"`
	Data     map[string]any `json:"data"`
}

func NewEvent(event storage.Event) Event {
	return Event{
		Id:       event.Id,
		Height:   event.Height,
		Time:     event.Time,
		Position: event.Position,
		Type:     string(event.Type),
		TxId:     event.TxId,
		Data:     event.Data,
	}
}
