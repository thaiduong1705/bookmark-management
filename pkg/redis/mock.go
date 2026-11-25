package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// InitMockRedis starts an in-memory miniredis and returns a client connected to it.
func InitMockRedis(t *testing.T) *redis.Client {
	t.Helper()
	mock := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client
}
