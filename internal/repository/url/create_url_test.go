package url

import (
	redisPkg "bookmark-management/pkg/redis"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestUrlShortenRepository_CreateUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		code        string
		url         string
		expr        int
		setupMock   func(t *testing.T) *redis.Client
		expectedErr error
	}{
		{
			name: "success",
			code: "abc",
			url:  "http://example.com",
			expr: 2,
			setupMock: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedErr: nil,
		},
		{
			name: "redis down",
			code: "abc",
			url:  "http://example.com",
			expr: 2,
			setupMock: func(t *testing.T) *redis.Client {
				rdb := redisPkg.InitMockRedis(t)
				rdb.Close()
				return rdb
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rdb := tc.setupMock(t)

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			testRepo := NewRepository(rdb)

			err := testRepo.CreateUrl(ctx, tc.code, tc.url, tc.expr)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)

			res, err := rdb.Get(ctx, tc.code).Result()
			assert.NoError(t, err)
			assert.Equal(t, tc.url, res)

			ttl, err := rdb.TTL(ctx, tc.code).Result()
			assert.NoError(t, err)
			assert.Greater(t, ttl, time.Duration(0))
			assert.LessOrEqual(t, ttl, time.Duration(tc.expr)*time.Minute)
		})
	}
}
