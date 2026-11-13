package service

import (
	"net/http/httptest"
	"testing"

	redisPkg "bookmark-management/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()

	t.Run("redis up", func(t *testing.T) {
		t.Parallel()
		rdb := redisPkg.InitMockRedis(t)
		defer rdb.Close()
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		assert.NoError(t, NewHealthCheck(rdb).Check(ctx))
	})

	t.Run("redis down", func(t *testing.T) {
		t.Parallel()
		rdb := redisPkg.InitMockRedis(t)
		rdb.Close()
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		assert.Error(t, NewHealthCheck(rdb).Check(ctx))
	})
}
