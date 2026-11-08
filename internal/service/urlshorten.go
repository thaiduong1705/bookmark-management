package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"shorten-url/internal/repository"
)

//go:generate mockery --name=UrlShortenService --filename urlshorten.go --outpkg urlshortenservice_svc_mocks
type UrlShortenService interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, url string, expr int) (string, error)
}

type urlShortenService struct {
	repo repository.UrlShortenRepository
}

func NewUrlShortenService(repo repository.UrlShortenRepository) UrlShortenService {
	return &urlShortenService{
		repo: repo,
	}
}

func (svc *urlShortenService) GetUrl(ctx context.Context, code string) (string, error) {
	return svc.repo.GetUrl(ctx, code)
}

func (svc *urlShortenService) CreateUrl(ctx context.Context, url string, expr int) (string, error) {
	code, err := generateRandomString(10)
	if err != nil {
		return "", err
	}

	if err = svc.repo.CreateUrl(ctx, code, url, expr); err != nil {
		return "", err
	}
	return code, nil
}

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
