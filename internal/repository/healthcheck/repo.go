package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repository checks the connectivity of the backing stores.
type Repository interface {
	Check(ctx context.Context) error
}

type repository struct {
	rdb    *redis.Client
	gormDB *gorm.DB
}

// NewRepository creates a healthcheck Repository that probes the given Redis client and database.
func NewRepository(rdb *redis.Client, gormDB *gorm.DB) Repository {
	return &repository{
		rdb:    rdb,
		gormDB: gormDB,
	}
}
