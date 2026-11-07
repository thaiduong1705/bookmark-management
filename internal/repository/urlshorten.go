package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UrlShortenRepository interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, code, url string, expr int) error
}

type urlShortenRepository struct {
	rdb *redis.Client
}

func NewUrlShortenRepository(rdb *redis.Client) UrlShortenRepository {
	return &urlShortenRepository{
		rdb: rdb,
	}
}

func (repo *urlShortenRepository) GetUrl(ctx context.Context, code string) (string, error) {
	return "", nil
}

func (repo *urlShortenRepository) CreateUrl(ctx context.Context, code, url string, expr int) error {
	if err := repo.rdb.Set(ctx, code, url, time.Duration(expr)*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}
