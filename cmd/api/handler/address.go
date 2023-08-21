package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	address storage.IAddress
}

func NewAddressHandler(address storage.IAddress) *AddressHandler {
	return &AddressHandler{
		address: address,
	}
}

type getAddressRequest struct {
	Hash string `param:"hash" validate:"required,address"`
}

func (handler *AddressHandler) Get(c echo.Context) error {
	req := new(getAddressRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	hash, err := DecodeAddress(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, address)
}

func (handler *AddressHandler) List(c echo.Context) error {
	req := new(limitOffsetPagination)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	address, err := handler.address.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}
	return returnArray(c, address)
}
