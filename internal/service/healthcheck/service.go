package healthcheck

import (
	"bookmark-management/internal/repository/healthcheck"
	"context"
)

// Service checks the health of the service dependencies.
type Service interface {
	Check(ctx context.Context) (string, string, error)
}

type service struct {
	repo        healthcheck.Repository
	serviceName string
	instanceId  string
}

// NewService creates a healthcheck Service that reports sName and insId when repo is healthy.
func NewService(sName string, insId string, repo healthcheck.Repository) Service {
	return &service{
		serviceName: sName,
		instanceId:  insId,
		repo:        repo,
	}
}
