package main

import (
	"shorten-url/internal/api"
	"shorten-url/pkg/redis"
)

// @title Bookmark Management API
// @version         1.0
// @description     API for managing bookmarks
// @BasePath /
func main() {
	cfg, err := api.NewConfig()
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
