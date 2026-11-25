package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Service checks the health of the service dependencies.
type Service interface {
	Check(ctx context.Context) (string, string, error)
}

type service struct {
	rdb         *redis.Client
	serviceName string
	instanceId  string
}

// NewService creates a HealthCheck backed by the given Redis client.
func NewService(rdb *redis.Client, sName string, insId string) Service {
	return &service{
		rdb:         rdb,
		serviceName: sName,
		instanceId:  insId,
	}
}
