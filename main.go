package main

import (
	"github.com/whaoa/public-api-bridge/internal/app"
	"github.com/whaoa/public-api-bridge/internal/config"
)

// @title Public API Bridge
// @version 1.0
// @description API bridge service for accessing public APIs (Twitter, Douban, TMDB, etc.)
// @termsOfService https://github.com/whaoa/public-api-bridge

// @contact.name API Support
// @contact.url https://github.com/whaoa/public-api-bridge/issues

// @BasePath /api
// @schemes http https

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	err = app.Serve(cfg)
	if err != nil {
		panic(err)
	}
}
