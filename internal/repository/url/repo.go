package url

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Repository stores and loads URLs by short code.
type Repository interface {
	GetUrl(ctx context.Context, code string) (string, error)
	CreateUrl(ctx context.Context, code, url string, expr int) error
}

type repository struct {
	rdb *redis.Client
}

// NewRepository creates a Redis-backed Repository.
func NewRepository(rdb *redis.Client) Repository {
	return &repository{
		rdb: rdb,
	}
}
