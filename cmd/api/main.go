package main

import (
	"bookmark-management/internal/api"
	"bookmark-management/pkg/redis"
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

	app := api.NewEngine(cfg, redisClient)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
