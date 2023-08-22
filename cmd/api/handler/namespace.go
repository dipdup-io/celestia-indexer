package handler

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type NamespaceHandler struct {
	namespace storage.INamespace
}

func NewNamespaceHandler(namespace storage.INamespace) *NamespaceHandler {
	return &NamespaceHandler{
		namespace: namespace,
	}
}

type getNamespaceRequest struct {
	Id string `param:"id" validate:"required,hexadecimal,len=56"`
}

func (handler *NamespaceHandler) Get(c echo.Context) error {
	req := new(getNamespaceRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	namespace, err := handler.namespace.ByNamespaceId(c.Request().Context(), namespaceId)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	response := make([]Namespace, len(namespace))
	for i := range namespace {
		response[i] = NewNamespace(namespace[i])
	}

	return returnArray(c, response)
}

type getNamespaceByHashRequest struct {
	Hash string `param:"hash" validate:"required,base64"`
}

func (handler *NamespaceHandler) GetByHash(c echo.Context) error {
	req := new(getNamespaceByHashRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}
	if len(hash) != 29 {
		return badRequestError(c, errors.Wrapf(errInvalidNamespaceLength, "got %d", len(hash)))
	}
	version := hash[0]
	namespaceId := hash[1:]

	namespace, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, version)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, NewNamespace(namespace))
}

type getNamespaceWithVersionRequest struct {
	Id      string `param:"id"      validate:"required,hexadecimal,len=56"`
	Version byte   `param:"version" validate:"required"`
}

func (handler *NamespaceHandler) GetWithVersion(c echo.Context) error {
	req := new(getNamespaceWithVersionRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	namespace, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, req.Version)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewNamespace(namespace))
}

func (handler *NamespaceHandler) List(c echo.Context) error {
	req := new(limitOffsetPagination)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	namespace, err := handler.namespace.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}
	response := make([]Namespace, len(namespace))
	for i := range namespace {
		response[i] = NewNamespace(*namespace[i])
	}
	return returnArray(c, response)
}
