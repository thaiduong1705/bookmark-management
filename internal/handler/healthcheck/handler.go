package healthcheck

import (
	service "bookmark-management/internal/service/healthcheck"

	"github.com/gin-gonic/gin"
)

// Handler serves the health check endpoint.
type Handler interface {
	Check(c *gin.Context)
}

type handler struct {
	srv service.Service
}

// NewHandler creates a HealthCheck handler.
func NewHandler(srv service.Service) Handler {
	return &handler{
		srv: srv,
	}
}
