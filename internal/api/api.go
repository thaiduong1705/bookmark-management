package api

import (
	"net/http"
	"shorten-url/internal/config"
	"shorten-url/internal/handler"
	"shorten-url/internal/service"

	"github.com/gin-gonic/gin"
)

type Engine interface {
	Start() error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *config.Config
}

func NewEngine(cfg *config.Config) Engine {
	e := &engine{
		app: gin.Default(),
		cfg: cfg,
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

	e.app.GET("/health-check", healthHandler.Check)
}
