package handler

import (
	"bookmark-management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck serves the health check endpoint.
type HealthCheck interface {
	Check(c *gin.Context)
}

type healthCheck struct {
	srv service.HealthCheck
}

// NewHealthCheck creates a HealthCheck handler.
func NewHealthCheck(srv service.HealthCheck) HealthCheck {
	return &healthCheck{
		srv: srv,
	}
}

// Check
// @Summary      Health check
// @Description  Checks that the service and its dependencies are healthy.
// @Tags         health
// @Produce      application/json
// @Success      200  {object}  map[string]string  "Service is healthy"
// @Failure      500  {object}  map[string]string  "Service is unhealthy"
// @Router       /health-check [get]
func (h *healthCheck) Check(c *gin.Context) {
	if err := h.srv.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unhealthy",
			"redis":   "down",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"redis":   "ok",
	})
}
