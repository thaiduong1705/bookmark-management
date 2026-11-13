package repository

import (
	"bookmark-management/internal/domain"
	redisPkg "bookmark-management/pkg/redis"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
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
			defer rdb.Close()

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			testRepo := NewUrlShortenRepository(rdb)

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

func TestUrlShortenRepository_GetUrl(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		code        string
		setupMock   func(t *testing.T) *redis.Client
		expectedErr error
		expectedRes string
	}{
		{
			name: "success",
			code: "abc",
			setupMock: func(t *testing.T) *redis.Client {
				rdb := redisPkg.InitMockRedis(t)
				rdb.Set(t.Context(), "abc", "http://example.com", 0)
				return rdb
			},
			expectedErr: nil,
			expectedRes: "http://example.com",
		},
		{
			name: "key not found",
			code: "empty",
			setupMock: func(t *testing.T) *redis.Client {
				rdb := redisPkg.InitMockRedis(t)
				rdb.Set(t.Context(), "abc", "http://example.com", 0)
				return rdb
			},
			expectedErr: domain.ErrNotFound,
			expectedRes: "",
		},
		{
			name: "key expired",
			code: "abc",
			setupMock: func(t *testing.T) *redis.Client {
				mr := miniredis.RunT(t)
				rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
				rdb.Set(t.Context(), "abc", "http://example.com", time.Minute)
				mr.FastForward(2 * time.Minute)
				return rdb
			},
			expectedErr: domain.ErrNotFound,
			expectedRes: "",
		},
		{
			name: "redis down",
			code: "abc",
			setupMock: func(t *testing.T) *redis.Client {
				rdb := redisPkg.InitMockRedis(t)
				rdb.Close()
				return rdb
			},
			expectedErr: redis.ErrClosed,
			expectedRes: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rdb := tc.setupMock(t)
			defer rdb.Close()
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			testRepo := NewUrlShortenRepository(rdb)
			res, err := testRepo.GetUrl(ctx, tc.code)
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}
