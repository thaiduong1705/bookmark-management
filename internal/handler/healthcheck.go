package handler

import (
	"net/http"
	"shorten-url/internal/service"

	"github.com/gin-gonic/gin"
)

type HealthCheck interface {
	Check(c *gin.Context)
}

type healthCheck struct {
	srv service.HealthCheck
}

func NewHealthCheck(srv service.HealthCheck) HealthCheck {
	return &healthCheck{
		srv: srv,
	}
}

func (h *healthCheck) Check(c *gin.Context) {
	serviceName, instanceID := h.srv.Check()
	c.JSON(http.StatusOK, gin.H{
		"message":      "ok",
		"service_name": serviceName,
		"instance_id":  instanceID,
	})
}
