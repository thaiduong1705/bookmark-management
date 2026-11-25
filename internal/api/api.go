package api

import (
	"bookmark-management/docs"
	healthcheckhandler "bookmark-management/internal/handler/healthcheck"
	urlshortenhandler "bookmark-management/internal/handler/url"
	userhandler "bookmark-management/internal/handler/user"
	healthcheckrepo "bookmark-management/internal/repository/healthcheck"
	urlshortenrepo "bookmark-management/internal/repository/url"
	userrepo "bookmark-management/internal/repository/user"
	healthcheckservice "bookmark-management/internal/service/healthcheck"
	urlshortenservice "bookmark-management/internal/service/url"
	userservice "bookmark-management/internal/service/user"
	"bookmark-management/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

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
	db  *gorm.DB
}

// NewEngine builds the Engine and registers all routes.
func NewEngine(cfg *Config, rdb *redis.Client, db *gorm.DB) Engine {
	app := gin.Default()
	app.ContextWithFallback = true

	e := &engine{
		app: app,
		cfg: cfg,
		rdb: rdb,
		db:  db,
	}
	e.initRoutes()

	return e
}

// Start runs the HTTP server on the configured port and blocks until it stops.
func (e *engine) Start() error {

	return e.app.Run(":" + e.cfg.Port)
}

// ServeHTTP lets the engine be used as an http.Handler, e.g. with httptest in tests.
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func (e *engine) initRoutes() {
	strGenUtil := utils.NewStringGenerator()
	passwordUtil := utils.NewPasswordUtil()

	healthRepo := healthcheckrepo.NewRepository(e.rdb, e.db)
	healthService := healthcheckservice.NewService(e.cfg.ServiceName, e.cfg.InstanceId, healthRepo)
	healthHandler := healthcheckhandler.NewHandler(healthService)

	urlShortenRepo := urlshortenrepo.NewRepository(e.rdb)
	urlShortenService := urlshortenservice.NewService(urlShortenRepo, strGenUtil)
	urlShortenHandler := urlshortenhandler.NewHandler(urlShortenService)

	userRepo := userrepo.NewRepository(e.db)
	userService := userservice.NewService(userRepo, passwordUtil)
	userHandler := userhandler.NewHandler(userService, passwordUtil)

	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", healthHandler.Check)

	v1 := e.app.Group("/v1")
	{
		v1.GET("/links/shorten/:code", urlShortenHandler.RedirectUrl)
		v1.POST("/links/shorten", urlShortenHandler.CreateUrl)
	}
	{
		v1.POST("/users/register", userHandler.Register)
	}
}
