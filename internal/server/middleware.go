package server

import (
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/whaoa/public-api-bridge/pkg/log"
)

func resolveRequestId(c *echo.Context) string {
	id := c.Request().Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = c.Response().Header().Get(echo.HeaderXRequestID)
	}
	return id
}

func shouldSkipLogRequest(c *echo.Context) bool {
	uri := c.Request().RequestURI
	if strings.HasPrefix(uri, "/.well-known/") {
		return true
	}
	if uri == "/swagger" || strings.HasPrefix(uri, "/swagger/") {
		return true
	}
	return false
}

func middlewareRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			request := c.Request()
			method := request.Method
			uri := request.RequestURI

			if shouldSkipLogRequest(c) {
				return next(c)
			}

			start := time.Now()
			ip := c.RealIP()

			err := next(c)

			id := resolveRequestId(c)
			latency := time.Now().Sub(start)
			_, status := echo.ResolveResponseStatus(c.Response(), err)

			c.Logger().InfoContext(
				request.Context(),
				"request",
				slog.String("rid", id),
				slog.String("method", method),
				slog.String("uri", uri),
				slog.Int("status", status),
				slog.Int64("latency", latency.Milliseconds()),
				slog.String("ip", ip),
				log.Err(err),
			)

			return err
		}
	}
}

func middlewareRequestLoggerMetadata() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if shouldSkipLogRequest(c) {
				return next(c)
			}

			request := c.Request()

			ctx := log.AddMetadataToContext(request.Context(), log.Metadata{
				"rid":    resolveRequestId(c),
				"time":   time.Now(),
				"method": request.Method,
				"uri":    request.RequestURI,
				"ip":     c.RealIP(),
			})

			c.SetRequest(request.WithContext(ctx))

			return next(c)
		}
	}
}
