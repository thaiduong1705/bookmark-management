package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg svc_mocks

// HealthCheck checks the health of the service dependencies.
type HealthCheck interface {
	Check(c *gin.Context) error
}

type healthCheck struct {
	rdb *redis.Client
}

// NewHealthCheck creates a HealthCheck backed by the given Redis client.
func NewHealthCheck(rdb *redis.Client) HealthCheck {
	return &healthCheck{rdb: rdb}
}

// Check pings Redis and returns an error if it is unreachable.
func (h *healthCheck) Check(c *gin.Context) error {
	if err := h.rdb.Ping(c).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}
