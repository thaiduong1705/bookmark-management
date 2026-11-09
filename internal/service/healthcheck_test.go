package service

import (
	"testing"

	redisPkg "bookmark-management/pkg/redis"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()

	t.Run("redis up", func(t *testing.T) {
		t.Parallel()
		rdb := redisPkg.InitMockRedis(t)
		defer rdb.Close()
		assert.NoError(t, NewHealthCheck(rdb).Check(t.Context()))
	})

	t.Run("redis down", func(t *testing.T) {
		t.Parallel()
		rdb := redisPkg.InitMockRedis(t)
		rdb.Close()
		assert.Error(t, NewHealthCheck(rdb).Check(t.Context()))
	})
}
