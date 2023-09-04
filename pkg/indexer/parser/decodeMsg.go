package parser

import (
	"github.com/celestiaorg/celestia-app/pkg/namespace"
	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func decodeMsg(b types.BlockData, msg cosmosTypes.Msg, position int) (storage.Message, error) {
	fullMsgType := reflect.TypeOf(msg).String()
	msgTypeName := fullMsgType[strings.LastIndex(fullMsgType, ".")+1:]
	msgType := storageTypes.MsgTypeUnknown
	if storageTypes.IsMsgType(msgTypeName) {
		msgType = storageTypes.MsgType(msgTypeName)
	}

	storageMsg := storage.Message{
		Height:   b.Height,
		Time:     b.Block.Time,
		Position: uint64(position),
		Type:     msgType,
		Data:     structs.Map(msg),
	}

	// Decode Namespaces
	if msgType == storageTypes.MsgTypePayForBlobs {
		payForBlobsMsg, ok := msg.(*appBlobTypes.MsgPayForBlobs)
		if !ok {
			return storage.Message{}, errors.Errorf("error on decoding %T", msg)
		}

		storageMsg.Namespace = make([]storage.Namespace, len(payForBlobsMsg.Namespaces))
		for nsI, ns := range payForBlobsMsg.Namespaces {
			if len(payForBlobsMsg.BlobSizes) < nsI {
				return storage.Message{}, errors.Errorf("blob sizes does not match with namespaces %d in msg on position %d", nsI, position)
			}

			appNS := namespace.Namespace{Version: ns[0], ID: ns[1:]}
			storageMsg.Namespace[nsI] = storage.Namespace{
				FirstHeight: b.Height,
				Version:     ns[0],
				NamespaceID: ns[1:],
				Size:        uint64(payForBlobsMsg.BlobSizes[nsI]),
				PfbCount:    1,
				Reserved:    appNS.IsReserved(),
			}
		}
	}

	return storageMsg, nil
}
