package api

import (
	"bookmark-management/internal/handler"
	"bookmark-management/internal/repository"
	"bookmark-management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "bookmark-management/docs"
)

// Engine is the HTTP server of the application.
type Engine interface {
	Start() error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
	rdb *redis.Client
}

// NewEngine builds the Engine and registers all routes.
func NewEngine(cfg *Config, rdb *redis.Client) Engine {
	app := gin.Default()
	app.ContextWithFallback = true

	e := &engine{
		app: app,
		cfg: cfg,
		rdb: rdb,
	}
	e.initRoutes()

	return e
}

func (e *engine) Start() error {

	return e.app.Run(":" + e.cfg.Port)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func (e *engine) initRoutes() {
	healthService := service.NewHealthCheck(e.rdb)
	healthHandler := handler.NewHealthCheck(healthService)

	urlShortenRepo := repository.NewUrlShortenRepository(e.rdb)
	urlShortenService := service.NewUrlShortenService(urlShortenRepo)
	urlShortenHandler := handler.NewUrlShortenHandler(urlShortenService)

	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", healthHandler.Check)
	e.app.GET("/v1/links/shorten/:code", urlShortenHandler.RedirectUrl)
	e.app.POST("/v1/links/shorten", urlShortenHandler.CreateUrl)
}
