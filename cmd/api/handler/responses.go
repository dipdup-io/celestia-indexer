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

// Address model info
// @Description Celestia address information
type Address struct {
	Id      uint64 `example:"321"                                             json:"id"           swaggertype:"integer"`
	Height  uint64 `example:"100"                                             json:"first_height" swaggertype:"integer"`
	Balance string `example:"10000000000"                                     json:"balance"      swaggertype:"string"`
	Hash    string `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" json:"hash"         swaggertype:"string"`
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
	Id                 uint64    `example:"321"                                                              json:"id"                   swaggertype:"integer"`
	Height             uint64    `example:"100"                                                              json:"height"               swaggertype:"integer"`
	Time               time.Time `example:"2023-07-04T03:10:57+00:00"                                        json:"time"                 swaggertype:"string"`
	VersionBlock       string    `example:"11"                                                               json:"version_block"        swaggertype:"string"`
	VersionApp         string    `example:"1"                                                                json:"version_app"          swaggertype:"string"`
	TxCount            uint64    `example:"12"                                                               json:"tx_count"             swaggertype:"integer"`
	EventsCount        uint64    `example:"18"                                                               json:"events_count"         swaggertype:"integer"`
	Hash               string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"hash"                 swaggertype:"string"`
	ParentHash         string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"parent_hash"          swaggertype:"string"`
	LastCommitHash     string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_commit_hash"     swaggertype:"string"`
	DataHash           string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"data_hash"            swaggertype:"string"`
	ValidatorsHash     string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"validators_hash"      swaggertype:"string"`
	NextValidatorsHash string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"next_validators_hash" swaggertype:"string"`
	ConsensusHash      string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"consensus_hash"       swaggertype:"string"`
	AppHash            string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"app_hash"             swaggertype:"string"`
	LastResultsHash    string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_results_hash"    swaggertype:"string"`
	EvidenceHash       string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"evidence_hash"        swaggertype:"string"`
	ProposerAddress    string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"proposer_address"     swaggertype:"string"`
}

func NewBlock(block storage.Block) Block {
	return Block{
		Id:                 block.Id,
		Height:             block.Height,
		Time:               block.Time,
		VersionBlock:       block.VersionBlock,
		VersionApp:         block.VersionApp,
		TxCount:            block.TxCount,
		EventsCount:        block.EventsCount,
		Hash:               hex.EncodeToString(block.Hash),
		ParentHash:         hex.EncodeToString(block.ParentHash),
		LastCommitHash:     hex.EncodeToString(block.LastCommitHash),
		DataHash:           hex.EncodeToString(block.DataHash),
		ValidatorsHash:     hex.EncodeToString(block.ValidatorsHash),
		NextValidatorsHash: hex.EncodeToString(block.NextValidatorsHash),
		ConsensusHash:      hex.EncodeToString(block.ConsensusHash),
		AppHash:            hex.EncodeToString(block.AppHash),
		LastResultsHash:    hex.EncodeToString(block.LastResultsHash),
		EvidenceHash:       hex.EncodeToString(block.EvidenceHash),
		ProposerAddress:    hex.EncodeToString(block.ProposerAddress),
	}
}

type Event struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"              swaggettype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"          swaggettype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggettype:"string"`
	Position uint64    `example:"1"                         format:"int64"     json:"position"        swaggettype:"integer"`
	TxId     uint64    `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggettype:"integer"`

	Type string `enums:"coin_received,coinbase,coin_spent,burn,mint,message,proposer_reward,rewards,commission,liveness,attestation_request,transfer,pay_for_blobs,redelegate,withdraw_rewards,withdraw_commission,create_validator,delegate,edit_validator,unbond,tx,unknown" example:"commission" format:"string" json:"type" swaggettype:"string"`

	Data map[string]any `json:"data"`
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
	Id       uint64    `example:"321"                       format:"int64"     json:"id"              swaggettype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"          swaggettype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggettype:"string"`
	Position uint64    `example:"2"                         format:"int64"     json:"position"        swaggettype:"integer"`
	TxId     uint64    `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggettype:"integer"`

	Type string `enums:"WithdrawValidatorCommission,WithdrawDelegatorReward,EditValidator,BeginRedelegate,CreateValidator,Delegate,Undelegate,Unjail,Send,CreateVestingAccount,CreatePeriodicVestingAccount,PayForBlobs" example:"CreatePeriodicVestingAccount" format:"string" json:"type" swaggettype:"string"`

	Data map[string]any `json:"data"`
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
	Id                 uint64    `example:"321"                       format:"int64"     json:"id"                   swaggettype:"integer"`
	Name               string    `example:"indexer"                   format:"string"    json:"name"                 swaggettype:"string"`
	LastHeight         uint64    `example:"100"                       format:"int64"     json:"last_height"          swaggettype:"integer"`
	LastTime           time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_time"            swaggettype:"string"`
	TotalTx            uint64    `example:"23456"                     format:"int64"     json:"total_tx"             swaggettype:"integer"`
	TotalAccounts      uint64    `example:"43"                        format:"int64"     json:"total_accounts"       swaggettype:"integer"`
	TotalNamespaces    uint64    `example:"312"                       format:"int64"     json:"total_namespaces"     swaggettype:"integer"`
	TotalNamespaceSize uint64    `example:"56789"                     format:"int64"     json:"total_namespace_size" swaggettype:"integer"`
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
	Id            uint64    `example:"321"                                                              format:"int64"     json:"id"                  swaggettype:"integer"`
	Height        uint64    `example:"100"                                                              format:"int64"     json:"height"              swaggettype:"integer"`
	Position      uint64    `example:"11"                                                               format:"int64"     json:"position"            swaggettype:"integer"`
	GasWanted     uint64    `example:"9348"                                                             format:"int64"     json:"gas_wanted"          swaggettype:"integer"`
	GasUsed       uint64    `example:"4253"                                                             format:"int64"     json:"gas_used"            swaggettype:"integer"`
	TimeoutHeight uint64    `example:"0"                                                                format:"int64"     json:"timeout_height"      swaggettype:"integer"`
	EventsCount   uint64    `example:"2"                                                                format:"int64"     json:"events_count"        swaggettype:"integer"`
	MessagesCount uint64    `example:"1"                                                                format:"int64"     json:"messages_count"      swaggettype:"integer"`
	Hash          string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"hash"                swaggettype:"string"`
	Fee           string    `example:"9348"                                                             format:"int64"     json:"fee"                 swaggettype:"string"`
	Error         string    `example:""                                                                 format:"string"    json:"error,omitempty"     swaggettype:"string"`
	Codespace     string    `example:"sdk"                                                              format:"string"    json:"codespace,omitempty" swaggettype:"string"`
	Memo          string    `example:"Transfer to private account"                                      format:"string"    json:"memo,omitempty"      swaggettype:"string"`
	Time          time.Time `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                swaggettype:"string"`

	Status string `enums:"success,failed" example:"success" format:"string" json:"status" swaggettype:"string"`
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
	ID          uint64 `example:"321"                                                      format:"integer" json:"id"           swaggertype:"integer"`
	Size        uint64 `example:"12345"                                                    format:"integer" json:"size"         swaggertype:"integer"`
	Version     byte   `examle:"1"                                                         format:"byte"    json:"version"      swaggertype:"integer"`
	NamespaceID string `example:"4723ce10b187716adfc55ff7e6d9179c226e6b5440b02577cca49d02" format:"binary"  json:"namespace_id" swaggertype:"string"`
	Hash        string `example:"U3dhZ2dlciByb2Nrcw=="                                     format:"base64"  json:"hash"         swaggertype:"string"`
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
