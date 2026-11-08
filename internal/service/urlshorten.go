package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"shorten-url/internal/domain"
	"shorten-url/internal/repository"
)

// UrlShortenService contains the URL shortening business logic.
//
//go:generate mockery --name=UrlShortenService --filename urlshorten.go --outpkg svc_mocks
type UrlShortenService interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, url string, expr int) (string, error)
}

type urlShortenService struct {
	repo repository.UrlShortenRepository
}

// NewUrlShortenService creates a UrlShortenService.
func NewUrlShortenService(repo repository.UrlShortenRepository) UrlShortenService {
	return &urlShortenService{
		repo: repo,
	}
}

// GetUrl returns the original URL for the given code.
func (svc *urlShortenService) GetUrl(ctx context.Context, code string) (string, error) {
	return svc.repo.GetUrl(ctx, code)
}

// CreateUrl generates a random code and stores the URL under it.
func (svc *urlShortenService) CreateUrl(ctx context.Context, url string, expr int) (string, error) {
	code, err := generateRandomString(domain.RandomLength)
	if err != nil {
		return "", err
	}

	if err = svc.repo.CreateUrl(ctx, code, url, expr); err != nil {
		return "", err
	}
	return code, nil
}

// generateRandomString returns a cryptographically random alphanumeric string.
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
