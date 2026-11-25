package healthcheck

import (
	redisPkg "bookmark-management/pkg/redis"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()

	const (
		serviceName = "bookmark_service"
		instanceId  = "instance_1"
	)

	t.Run("redis up", func(t *testing.T) {
		t.Parallel()
		rdb := redisPkg.InitMockRedis(t)

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		name, id, err := NewService(rdb, serviceName, instanceId).Check(ctx)

		assert.NoError(t, err)
		assert.Equal(t, serviceName, name)
		assert.Equal(t, instanceId, id)
	})

	t.Run("redis down", func(t *testing.T) {
		t.Parallel()
		rdb := redisPkg.InitMockRedis(t)
		rdb.Close()

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		name, id, err := NewService(rdb, serviceName, instanceId).Check(ctx)

		assert.Error(t, err)
		assert.Empty(t, name)
		assert.Empty(t, id)
	})
}
