package handler

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
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
//
//	@Summary		Get transaction by hash
//	@Description	Get transaction by hash
//	@Tags			transactions
//	@ID				get-transaction
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{object}	responses.Tx
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash} [get]
func (handler *TxHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
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

	return c.JSON(http.StatusOK, responses.NewTx(tx))
}

// List godoc
//
//	@Summary		List transactions info
//	@Description	List transactions info
//	@Tags			transactions
//	@ID				list-transactions
//	@Param			limit		query	integer	false	"Count of requested entities"			mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"								mininum(1)
//	@Param			sort		query	string	false	"Sort order"							Enums(asc, desc)
//	@Param			status		query	string	false	"Comma-separated status list"			Enums(success, failed)
//	@Param			msg_type	query	string	false	"Comma-separated message types list"	Enums(WithdrawValidatorCommission,WithdrawDelegatorReward,EditValidator,BeginRedelegate,CreateValidator,Delegate,Undelegate,Unjail,Send,CreateVestingAccount,CreatePeriodicVestingAccount,PayForBlobs)
//	@Param			from		query	integer	false	"Time from in unix timestamp"			mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"				mininum(1)
//	@Param			height		query	integer	false	"Block number"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx [get]
func (handler *TxHandler) List(c echo.Context) error {
	req, err := bindAndValidate[txListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.TxFilter{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
		Sort:   pgSort(req.Sort),
		Status: req.Status,
		Height: req.Height,
	}
	if req.From > 0 {
		fltrs.TimeFrom = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		fltrs.TimeTo = time.Unix(req.To, 0).UTC()
	}

	txs, err := handler.tx.Filter(c.Request().Context(), fltrs)
	if err := handleError(c, err, handler.tx); err != nil {
		return err
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

// GetEvents godoc
//
//	@Summary		Get transaction events
//	@Description	Get transaction events
//	@Tags			transactions
//	@ID				get-transaction-events
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{array}		responses.Event
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/events [get]
func (handler *TxHandler) GetEvents(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
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
	response := make([]responses.Event, len(events))
	for i := range events {
		response[i] = responses.NewEvent(events[i])
	}
	return returnArray(c, response)
}

// GetMessages godoc
//
//	@Summary		Get transaction messages
//	@Description	Get transaction messages
//	@Tags			transactions
//	@ID				get-transaction-messages
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{array}		responses.Message
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/messages [get]
func (handler *TxHandler) GetMessages(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
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
	response := make([]responses.Message, len(messages))
	for i := range messages {
		response[i] = responses.NewMessage(messages[i])
	}
	return returnArray(c, response)
}
