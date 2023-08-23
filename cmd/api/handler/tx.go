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
	Hash string `param:"hash" validate:"required,hexadecimal,len=64"`
}

// Get godoc
// @Summary Get transaction by hash
// @Description Get transaction by hash
// @Tags transactions
// @ID get-transaction
// @Param hash path string true "Transaction hash in hexadecimal" minlength(64) maxlength(64)
// @Produce  json
// @Success 200 {object} Tx
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/tx/{hash} [get]
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

	return c.JSON(http.StatusOK, NewTx(tx))
}

// List godoc
// @Summary List transactions info
// @Description List transactions info
// @Tags transactions
// @ID list-transactions
// @Param limit  query integer false "Count of requested entities" mininum(1) maximum(100)
// @Param offset query integer false "Offset" mininum(1)
// @Param sort   query string  false "Sort order" Enums(asc, desc)
// @Produce json
// @Success 200 {array} Tx
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/tx [get]
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
	response := make([]Tx, len(txs))
	for i := range txs {
		response[i] = NewTx(*txs[i])
	}
	return returnArray(c, response)
}

// GetEvents godoc
// @Summary Get transaction events
// @Description Get transaction events
// @Tags transactions
// @ID get-transaction-events
// @Param hash path string true "Transaction hash in hexadecimal" minlength(64) maxlength(64)
// @Produce json
// @Success 200 {array} Event
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/tx/{hash}/events [get]
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
	response := make([]Event, len(events))
	for i := range events {
		response[i] = NewEvent(events[i])
	}
	return returnArray(c, response)
}

// GetMessages godoc
// @Summary Get transaction messages
// @Description Get transaction messages
// @Tags transactions
// @ID get-transaction-messages
// @Param hash path string true "Transaction hash in hexadecimal" minlength(64) maxlength(64)
// @Produce json
// @Success 200 {array} Message
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/tx/{hash}/messages [get]
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
	response := make([]Message, len(messages))
	for i := range messages {
		response[i] = NewMessage(messages[i])
	}
	return returnArray(c, response)
}
