package repository

import (
	"bookmark-management/internal/domain"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// UrlShortenRepository stores and loads URLs by short code.
//
//go:generate mockery --name=UrlShortenRepository --filename urlshorten.go --outpkg repo_mocks
type UrlShortenRepository interface {
	GetUrl(c *gin.Context, code string) (string, error)
	CreateUrl(c *gin.Context, code, url string, expr int) error
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
func (repo *urlShortenRepository) GetUrl(c *gin.Context, code string) (string, error) {
	url, err := repo.rdb.Get(c, code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrNotFound
		}
		return "", err
	}
	return url, nil
}

// CreateUrl stores url under code; expr is in minutes (0 = never expires).
func (repo *urlShortenRepository) CreateUrl(c *gin.Context, code, url string, expr int) error {
	if err := repo.rdb.Set(c, code, url, time.Duration(expr)*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}
