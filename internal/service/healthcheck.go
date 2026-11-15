package service

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg svc_mocks

// HealthCheck checks the health of the service dependencies.
type HealthCheck interface {
	Check(ctx context.Context) (string, string, error)
}

type healthCheck struct {
	rdb         *redis.Client
	serviceName string
	instanceId  string
}

// NewHealthCheck creates a HealthCheck backed by the given Redis client.
func NewHealthCheck(rdb *redis.Client, sName string, insId string) HealthCheck {
	return &healthCheck{
		rdb:         rdb,
		serviceName: sName,
		instanceId:  insId,
	}
}

// Check pings Redis and returns an error if it is unreachable.
func (h *healthCheck) Check(ctx context.Context) (string, string, error) {
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		return "", "", fmt.Errorf("redis ping: %w", err)
	}
	return h.serviceName, h.instanceId, nil
}
