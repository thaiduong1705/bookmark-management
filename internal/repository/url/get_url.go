package url

import (
	"bookmark-management/internal/domain"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// GetUrl returns the URL for code, or domain.ErrNotFound if missing.
func (repo *repository) GetUrl(ctx context.Context, code string) (string, error) {
	url, err := repo.rdb.Get(ctx, code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrNotFound
		}
		return "", err
	}
	return url, nil
}
