package healthcheck

import (
	redisPkg "bookmark-management/pkg/redis"
	"bookmark-management/pkg/sqldb"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Check(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		setup       func(t *testing.T) (*redis.Client, *gorm.DB)
		expectedErr error
		expectedMsg string
	}{
		{
			name: "success",
			setup: func(t *testing.T) (*redis.Client, *gorm.DB) {
				return redisPkg.InitMockRedis(t), sqldb.InitMockDatabase(t)
			},
		},
		{
			name: "redis down",
			setup: func(t *testing.T) (*redis.Client, *gorm.DB) {
				rdb := redisPkg.InitMockRedis(t)
				rdb.Close()
				return rdb, sqldb.InitMockDatabase(t)
			},
			expectedErr: redis.ErrClosed,
			expectedMsg: "ping redis",
		},
		{
			name: "cannot get sql db",
			setup: func(t *testing.T) (*redis.Client, *gorm.DB) {
				db := sqldb.InitMockDatabase(t)
				db.ConnPool = nil
				db.Statement.ConnPool = nil
				return redisPkg.InitMockRedis(t), db
			},
			expectedErr: gorm.ErrInvalidDB,
			expectedMsg: "sql db",
		},
		{
			name: "database down",
			setup: func(t *testing.T) (*redis.Client, *gorm.DB) {
				db := sqldb.InitMockDatabase(t)
				sqlDB, err := db.DB()
				assert.NoError(t, err)
				assert.NoError(t, sqlDB.Close())
				return redisPkg.InitMockRedis(t), db
			},
			expectedMsg: "sql db",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rdb, db := tc.setup(t)

			err := NewRepository(rdb, db).Check(t.Context())

			if tc.expectedMsg == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tc.expectedMsg)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
