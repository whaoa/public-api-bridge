package main

import (
	"github.com/whaoa/public-api-bridge/internal/app"
	"github.com/whaoa/public-api-bridge/internal/config"
)

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
