package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	errInvalidNamespaceLength = errors.New("invalid namespace: should be 29 bytes length")
)

type NoRows interface {
	IsNoRows(err error) bool
}

type Error struct {
	Message string `json:"message"`
}

func badRequestError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, Error{
		Message: err.Error(),
	})
}

func internalServerError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, Error{
		Message: err.Error(),
	})
}

func handleError(c echo.Context, err error, noRows NoRows) error {
	if err == nil {
		return nil
	}
	if noRows.IsNoRows(err) {
		return c.NoContent(http.StatusNoContent)
	}
	return internalServerError(c, err)
}

func returnArray[T any](c echo.Context, arr []T) error {
	if arr == nil {
		return c.JSON(http.StatusOK, []any{})
	}

	return c.JSON(http.StatusOK, arr)
}
