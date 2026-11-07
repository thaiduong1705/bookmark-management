package api

import (
	"net/http"
	"shorten-url/internal/config"
	"shorten-url/internal/handler"
	"shorten-url/internal/repository"
	"shorten-url/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Engine interface {
	Start() error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *config.Config
	rdb *redis.Client
}

func NewEngine(cfg *config.Config, rdb *redis.Client) Engine {
	e := &engine{
		app: gin.Default(),
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
	healthService := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID)
	healthHandler := handler.NewHealthCheck(healthService)

	urlShortenRepo := repository.NewUrlShortenRepository(e.rdb)
	urlShortenService := service.NewUrlShortenService(urlShortenRepo)
	urlShortenHandler := handler.NewUrlShortenHandler(urlShortenService)

	e.app.GET("/health-check", healthHandler.Check)
	e.app.GET("/v1/links/shorten", urlShortenHandler.GetUrl)
	e.app.POST("/v1/links/shorten", urlShortenHandler.CreateUrl)
}
