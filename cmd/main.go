package main

import (
	"shorten-url/internal/api"
	"shorten-url/internal/config"
	"shorten-url/pkg/redis"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	app := api.NewEngine(cfg, redisClient)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
