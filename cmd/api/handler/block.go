package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type BlockHandler struct {
	block  storage.IBlock
	events storage.IEvent
}

func NewBlockHandler(block storage.IBlock, events storage.IEvent) *BlockHandler {
	return &BlockHandler{
		block:  block,
		events: events,
	}
}

type getBlockRequest struct {
	Height uint64 `param:"height" validate:"required,min=1"`
}

func (handler *BlockHandler) Get(c echo.Context) error {
	req := new(getBlockRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	block, err := handler.block.ByHeight(c.Request().Context(), req.Height)
	if err := handleError(c, err, handler.block); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewBlock(block))
}

func (handler *BlockHandler) List(c echo.Context) error {
	req := new(limitOffsetPagination)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blocks, err := handler.block.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.block); err != nil {
		return err
	}

	response := make([]Block, len(blocks))
	for i := range blocks {
		response[i] = NewBlock(*blocks[i])
	}

	return returnArray(c, response)
}

func (handler *BlockHandler) GetEvents(c echo.Context) error {
	req := new(getBlockRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	events, err := handler.events.ByBlock(c.Request().Context(), req.Height)
	if err := handleError(c, err, handler.events); err != nil {
		return err
	}

	response := make([]Event, len(events))
	for i := range events {
		response[i] = NewEvent(events[i])
	}

	return returnArray(c, response)
}
