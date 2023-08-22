package handler

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

func returnArray[T any](c echo.Context, arr []T) error {
	if arr == nil {
		return c.JSON(http.StatusOK, []any{})
	}

	return c.JSON(http.StatusOK, arr)
}

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
	TxId     uint64         `json:"tx_id,omitempty"`
	Data     map[string]any `json:"data"`
}

func NewEvent(event storage.Event) Event {
	result := Event{
		Id:       event.Id,
		Height:   event.Height,
		Time:     event.Time,
		Position: event.Position,
		Type:     string(event.Type),
		Data:     event.Data,
	}

	if event.TxId != nil {
		result.TxId = *event.TxId
	}

	return result
}

type Message struct {
	Id       uint64         `json:"id"`
	Height   uint64         `json:"height"`
	Time     time.Time      `json:"time"`
	Position uint64         `json:"position"`
	Type     string         `json:"type"`
	TxId     uint64         `json:"tx_id,omitempty"`
	Data     map[string]any `json:"data"`
}

func NewMessage(msg storage.Message) Message {
	return Message{
		Id:       msg.Id,
		Height:   msg.Height,
		Time:     msg.Time,
		Position: msg.Position,
		Type:     string(msg.Type),
		TxId:     msg.TxId,
		Data:     msg.Data,
	}
}

type State struct {
	Id                 uint64    `json:"id"`
	Name               string    `json:"name"`
	LastHeight         uint64    `json:"last_height"`
	LastTime           time.Time `json:"last_time"`
	TotalTx            uint64    `json:"total_tx"`
	TotalAccounts      uint64    `json:"total_accounts"`
	TotalNamespaces    uint64    `json:"total_namespaces"`
	TotalNamespaceSize uint64    `json:"total_namespace_size"`
}

func NewState(state storage.State) State {
	return State{

		Id:                 state.ID,
		Name:               state.Name,
		LastHeight:         state.LastHeight,
		LastTime:           state.LastTime,
		TotalTx:            state.TotalTx,
		TotalAccounts:      state.TotalAccounts,
		TotalNamespaces:    state.TotalNamespaces,
		TotalNamespaceSize: state.TotalNamespaceSize,
	}
}

type Tx struct {
	Id            uint64    `json:"id"`
	Height        uint64    `json:"height"`
	Position      uint64    `json:"position"`
	GasWanted     uint64    `json:"gas_wanted"`
	GasUsed       uint64    `json:"gas_used"`
	TimeoutHeight uint64    `json:"timeout_height"`
	EventsCount   uint64    `json:"events_count"`
	MessagesCount uint64    `json:"messages_count"`
	Hash          string    `json:"hash"`
	Fee           string    `json:"fee"`
	Status        string    `json:"status"`
	Error         string    `json:"error,omitempty"`
	Codespace     string    `json:"codespace,omitempty"`
	Memo          string    `json:"memo,omitempty"`
	Time          time.Time `json:"time"`
}

func NewTx(tx storage.Tx) Tx {
	return Tx{
		Id:            tx.Id,
		Height:        tx.Height,
		Time:          tx.Time,
		Position:      tx.Position,
		GasWanted:     tx.GasWanted,
		GasUsed:       tx.GasUsed,
		TimeoutHeight: tx.TimeoutHeight,
		EventsCount:   tx.EventsCount,
		MessagesCount: tx.MessagesCount,
		Fee:           tx.Fee.String(),
		Status:        string(tx.Status),
		Error:         tx.Error,
		Codespace:     tx.Codespace,
		Hash:          hex.EncodeToString(tx.Hash),
		Memo:          tx.Memo,
	}
}

type Namespace struct {
	ID          uint64 `json:"id"`
	Size        uint64 `json:"size"`
	Version     byte   `json:"version"`
	NamespaceID string `json:"namespace_id"`
	Hash        string `json:"hash"`
}

func NewNamespace(ns storage.Namespace) Namespace {
	return Namespace{
		ID:          ns.ID,
		Size:        ns.Size,
		Version:     ns.Version,
		NamespaceID: hex.EncodeToString(ns.NamespaceID),
		Hash:        base64.URLEncoding.EncodeToString(append([]byte{ns.Version}, ns.NamespaceID...)),
	}
}
