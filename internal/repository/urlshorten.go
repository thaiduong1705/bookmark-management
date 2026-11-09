package repository

import (
	"bookmark-management/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// UrlShortenRepository stores and loads URLs by short code.
//
//go:generate mockery --name=UrlShortenRepository --filename urlshorten.go --outpkg repo_mocks
type UrlShortenRepository interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, code, url string, expr int) error
}

type urlShortenRepository struct {
	rdb *redis.Client
}

// NewUrlShortenRepository creates a Redis-backed UrlShortenRepository.
func NewUrlShortenRepository(rdb *redis.Client) UrlShortenRepository {
	return &urlShortenRepository{
		rdb: rdb,
	}
}

// GetUrl returns the URL for code, or domain.ErrNotFound if missing.
func (repo *urlShortenRepository) GetUrl(ctx context.Context, code string) (string, error) {
	url, err := repo.rdb.Get(ctx, code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrNotFound
		}
		return "", err
	}
	return url, nil
}

// CreateUrl stores url under code; expr is in minutes (0 = never expires).
func (repo *urlShortenRepository) CreateUrl(ctx context.Context, code, url string, expr int) error {
	if err := repo.rdb.Set(ctx, code, url, time.Duration(expr)*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}
