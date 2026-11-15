package service

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/repository"
	"bookmark-management/pkg/utils"
	"context"
	"errors"
)

//go:generate mockery --name=UrlShortenService --filename urlshorten.go --outpkg svc_mocks

// UrlShortenService contains the URL shortening business logic.
type UrlShortenService interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, url string, expr int) (string, error)
}

type urlShortenService struct {
	repo       repository.UrlShortenRepository
	strGenUtil utils.StringGenerator
}

// NewUrlShortenService creates a UrlShortenService.
func NewUrlShortenService(repo repository.UrlShortenRepository, strGenUtil utils.StringGenerator) UrlShortenService {
	return &urlShortenService{
		repo:       repo,
		strGenUtil: strGenUtil,
	}
}

// GetUrl returns the original URL for the given code.
func (svc *urlShortenService) GetUrl(ctx context.Context, code string) (string, error) {
	return svc.repo.GetUrl(ctx, code)
}

// CreateUrl generates a random code and stores the URL under it.
func (svc *urlShortenService) CreateUrl(ctx context.Context, url string, expr int) (string, error) {
	code, err := svc.strGenUtil.Generate(domain.RandomLength)
	if err != nil {
		return "", err
	}

	// check if the code exists, recursive it
	_, err = svc.repo.GetUrl(ctx, code)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", err
	} else if err == nil {
		return svc.CreateUrl(ctx, url, expr)
	}

	if err = svc.repo.CreateUrl(ctx, code, url, expr); err != nil {
		return "", err
	}
	return code, nil
}
