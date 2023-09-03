package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

type NamespaceMessage struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"       swaggertype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"   swaggertype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"     swaggertype:"string"`
	Position uint64    `example:"2"                         format:"int64"     json:"position" swaggertype:"integer"`

	Type string `enums:"WithdrawValidatorCommission,WithdrawDelegatorReward,EditValidator,BeginRedelegate,CreateValidator,Delegate,Undelegate,Unjail,Send,CreateVestingAccount,CreatePeriodicVestingAccount,PayForBlobs" example:"CreatePeriodicVestingAccount" format:"string" json:"type" swaggertype:"string"`

	Data map[string]any `json:"data"`
	Tx   Tx             `json:"tx"`
}

func NewNamespaceMessage(msg storage.NamespaceMessage) (NamespaceMessage, error) {
	if msg.Message == nil {
		return NamespaceMessage{}, errors.New("nil message in namespace message constructor")
	}
	if msg.Tx == nil {
		return NamespaceMessage{}, errors.New("nil tx in namespace message constructor")
	}

	return NamespaceMessage{
		Id:       msg.Message.Id,
		Height:   uint64(msg.Message.Height),
		Time:     msg.Message.Time,
		Position: msg.Message.Position,
		Type:     string(msg.Message.Type),
		Data:     msg.Message.Data,
		Tx:       NewTx(*msg.Tx),
	}, nil
}
