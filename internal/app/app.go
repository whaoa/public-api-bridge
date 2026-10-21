package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samber/lo"
	"github.com/whaoa/public-api-bridge/internal/config"
	"github.com/whaoa/public-api-bridge/internal/douban"
	doubanOfficial "github.com/whaoa/public-api-bridge/internal/douban/provider/official"
	"github.com/whaoa/public-api-bridge/internal/server"
	"github.com/whaoa/public-api-bridge/internal/tmdb"
	tmdbOfficial "github.com/whaoa/public-api-bridge/internal/tmdb/provider/official"
	"github.com/whaoa/public-api-bridge/internal/twitter"
	"github.com/whaoa/public-api-bridge/internal/twitter/provider/fx"
	"github.com/whaoa/public-api-bridge/internal/twitter/provider/teapot"
	"github.com/whaoa/public-api-bridge/internal/twitter/provider/vx"
	"github.com/whaoa/public-api-bridge/pkg/log"
	"github.com/whaoa/public-api-bridge/pkg/request"
)

func Serve(cfg *config.Config) error {
	logger := log.New(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := request.WithDefaultHeaders()
	client.SetTimeout(time.Second * 5)

	s := server.New(cfg.ServerPort, logger.With(log.Scope("http")))

	api := s.Group("/api")

	// Douban service
	doubanService := douban.NewService(doubanOfficial.NewProvider(client))
	douban.NewHandler(doubanService).RegisterRoutes(api)

	// TMDB service
	if len(cfg.TMDBKeys) == 0 {
		logger.Warn("tmdb api keys are not configured, tmdb endpoints will be unavailable")
	}
	tmdbService := tmdb.NewService(lo.Map(cfg.TMDBKeys, func(key string, _ int) tmdb.Provider {
		return tmdbOfficial.NewProvider(client, key)
	})...)
	tmdb.NewHandler(tmdbService).RegisterRoutes(api)

	// Twitter service
	twitterService := twitter.NewService(
		fx.NewProvider(client),
		vx.NewProvider(client),
		teapot.NewProvider(client),
	)
	twitter.NewHandler(twitterService).RegisterRoutes(api)

	return s.Serve(ctx, func(s *server.Server) error {
		s.Logger.Info("server start", slog.Int("port", cfg.ServerPort))
		return nil
	})
}
