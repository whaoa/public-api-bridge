package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

func JSON(c *echo.Context, payload any, err error) error {
	if e, ok := errors.AsType[*errors.ReasonError](err); ok {
		return c.JSON(e.Code, e)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Internal("internal server error"))
	}
	if payload == nil {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, payload)
}
