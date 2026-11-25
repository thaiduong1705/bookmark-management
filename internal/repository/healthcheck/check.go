package healthcheck

import (
	"context"
	"fmt"
)

// Check pings Redis and the SQL database, returning an error if either is unreachable.
func (r *repository) Check(ctx context.Context) error {
	if err := r.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	sqlDB, err := r.gormDB.DB()
	if err != nil {
		return fmt.Errorf("sql db: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("sql db: %w", err)
	}
	return nil
}
