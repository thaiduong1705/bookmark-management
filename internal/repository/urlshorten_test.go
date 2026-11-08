package repository

import (
	"context"
	redisPkg "shorten-url/pkg/redis"
	"testing"

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
		expectedErr error
	}{
		{
			name:        "success",
			code:        "abc",
			url:         "http://example.com",
			expr:        1,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			defer mockRedis.Close()

			ctx := t.Context()
			testRepo := NewUrlShortenRepository(mockRedis)
			err := testRepo.CreateUrl(ctx, tc.code, tc.url, tc.expr)
			assert.ErrorIs(t, err, tc.expectedErr)
			res, err := mockRedis.Get(ctx, tc.code).Result()
			assert.NoError(t, err)
			assert.Equal(t, tc.url, res)
		})
	}
}

func TestUrlShortenRepository_GetUrl(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		code        string
		url         string
		expr        int
		setup       func(ctx context.Context, m *redis.Client)
		expectedErr error
		expectedRes string
	}{
		{
			name: "success",
			code: "abc",
			url:  "http://example.com",
			expr: 0,
			setup: func(ctx context.Context, m *redis.Client) {
				m.Set(ctx, "abc", "http://example.com", 0)
			},
			expectedErr: nil,
			expectedRes: "http://example.com",
		},
		{
			name: "key not found",
			code: "empty",
			url:  "empty value",
			expr: 0,
			setup: func(ctx context.Context, m *redis.Client) {
				m.Set(ctx, "abc", "http://example.com", 0)
			},
			expectedErr: redis.Nil,
			expectedRes: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			defer mockRedis.Close()
			ctx := t.Context()

			tc.setup(ctx, mockRedis)
			testRepo := NewUrlShortenRepository(mockRedis)
			res, err := testRepo.GetUrl(ctx, tc.code)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}
