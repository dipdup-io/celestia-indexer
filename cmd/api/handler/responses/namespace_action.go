package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

type NamespaceAction struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"       swaggettype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"   swaggettype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"     swaggettype:"string"`
	Position uint64    `example:"2"                         format:"int64"     json:"position" swaggettype:"integer"`

	Type string `enums:"WithdrawValidatorCommission,WithdrawDelegatorReward,EditValidator,BeginRedelegate,CreateValidator,Delegate,Undelegate,Unjail,Send,CreateVestingAccount,CreatePeriodicVestingAccount,PayForBlobs" example:"CreatePeriodicVestingAccount" format:"string" json:"type" swaggettype:"string"`

	Data map[string]any `json:"data"`
	Tx   Tx             `json:"tx"`
}

func NewNamespaceAction(action storage.NamespaceAction) (NamespaceAction, error) {
	if action.Message == nil {
		return NamespaceAction{}, errors.New("nil message in nmaespace action constructor")
	}
	if action.Tx == nil {
		return NamespaceAction{}, errors.New("nil tx in nmaespace action constructor")
	}

	return NamespaceAction{
		Id:       action.Message.Id,
		Height:   action.Message.Height,
		Time:     action.Message.Time,
		Position: action.Message.Position,
		Type:     string(action.Message.Type),
		Data:     action.Message.Data,
		Tx:       NewTx(*action.Tx),
	}, nil
}
