package url

import (
	urlShortenRepo "bookmark-management/internal/repository/url"
	"bookmark-management/pkg/utils"
	"context"
)

// Service contains the URL shortening business logic.
type Service interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, url string, expr int) (string, error)
}

type service struct {
	repo       urlShortenRepo.Repository
	strGenUtil utils.StringGenerator
}

// NewService creates a Service.
func NewService(repo urlShortenRepo.Repository, strGenUtil utils.StringGenerator) Service {
	return &service{
		repo:       repo,
		strGenUtil: strGenUtil,
	}
}
