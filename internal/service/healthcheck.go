package service

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// HealthCheck checks the health of the service dependencies.
//
//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg svc_mocks
type HealthCheck interface {
	Check(ctx context.Context) error
}

type healthCheck struct {
	rdb *redis.Client
}

// NewHealthCheck creates a HealthCheck backed by the given Redis client.
func NewHealthCheck(rdb *redis.Client) HealthCheck {
	return &healthCheck{rdb: rdb}
}

// Check pings Redis and returns an error if it is unreachable.
func (h *healthCheck) Check(ctx context.Context) error {
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}
