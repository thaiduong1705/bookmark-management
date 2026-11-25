package main

import (
	"bookmark-management/internal/api"
	"bookmark-management/internal/domain"
	"bookmark-management/pkg/logger"
	"bookmark-management/pkg/redis"
	"bookmark-management/pkg/sqldb"
)

// @title Bookmark Management API
// @version         1.0
// @description     API for managing bookmarks
// @BasePath /
func main() {
	cfg, err := api.NewConfig("APP_")
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("REDIS_")
	if err != nil {
		panic(err)
	}

	db, err := sqldb.NewClient("DB_")
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		panic(err)
	}

	logger.SetLogLevel(cfg.LogLevel)

	app := api.NewEngine(cfg, redisClient, db)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
