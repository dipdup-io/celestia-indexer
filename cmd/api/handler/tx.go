package handler

import (
	"encoding/hex"
	"net/http"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type TxHandler struct {
	tx       storage.ITx
	events   storage.IEvent
	messages storage.IMessage
}

func NewTxHandler(tx storage.ITx, events storage.IEvent, messages storage.IMessage) *TxHandler {
	return &TxHandler{
		tx:       tx,
		events:   events,
		messages: messages,
	}
}

type getTxRequest struct {
	Hash string `param:"hash" validate:"required,len=hexadecimal,len=64"`
}

func (handler *TxHandler) Get(c echo.Context) error {
	req := new(getTxRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tx)
}

func (handler *TxHandler) List(c echo.Context) error {
	req := new(limitOffsetPagination)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	txs, err := handler.tx.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}
	return returnArray(c, txs)
}

func (handler *TxHandler) GetEvents(c echo.Context) error {
	req := new(getTxRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}

	events, err := handler.events.ByTxId(c.Request().Context(), tx.Id)
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}
	return returnArray(c, events)
}

func (handler *TxHandler) GetMessages(c echo.Context) error {
	req := new(getTxRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}

	messages, err := handler.messages.ByTxId(c.Request().Context(), tx.Id)
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}
	return returnArray(c, messages)
}
