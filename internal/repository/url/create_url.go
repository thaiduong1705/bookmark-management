package url

import (
	"context"
	"time"
)

// CreateUrl stores url under code; expr is in seconds (0 = never expires).
func (repo *repository) CreateUrl(ctx context.Context, code, url string, expr int) error {
	if err := repo.rdb.Set(ctx, code, url, time.Duration(expr)*time.Second).Err(); err != nil {
		return err
	}
	return nil
}
