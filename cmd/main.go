package main

import (
	"shorten-url/internal/api"
	"shorten-url/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
