package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Server struct {
	*echo.Echo
	port   int
	logger *slog.Logger
}

func (s *Server) Serve(ctx context.Context, beforeServe func(s *Server) error) error {
	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", s.port),
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: time.Second * 5,
		BeforeServeFunc: func(sc *http.Server) error {
			if beforeServe == nil {
				return nil
			}
			return beforeServe(s)
		},
	}
	return sc.Start(ctx, s.Echo)
}

func New(port int, logger *slog.Logger) *Server {
	engine := echo.NewWithConfig(echo.Config{
		Logger:           logger,
		HTTPErrorHandler: echo.DefaultHTTPErrorHandler(false),
	})

	engine.Use(middleware.Recover())

	return &Server{
		Echo:   engine,
		port:   port,
		logger: logger,
	}
}
