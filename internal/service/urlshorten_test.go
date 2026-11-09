package service

import (
	"bookmark-management/internal/domain"
	repo_mocks "bookmark-management/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errRepo = errors.New("repo error")

func TestUrlShortenService_CreateUrl(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		url         string
		expr        int
		setupMock   func(t *testing.T) *repo_mocks.UrlShortenRepository
		expectedErr error
	}{
		{
			name: "success",
			url:  "http://test.com",
			expr: 1,
			setupMock: func(t *testing.T) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("CreateUrl",
					mock.Anything,
					mock.MatchedBy(func(code string) bool { return len(code) == domain.RandomLength }),
					"http://test.com",
					1,
				).Return(nil)
				return mockRepo
			},
			expectedErr: nil,
		},
		{
			name: "repo error",
			url:  "http://test.com",
			expr: 1,
			setupMock: func(t *testing.T) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("CreateUrl",
					mock.Anything,
					mock.MatchedBy(func(code string) bool { return len(code) == domain.RandomLength }),
					"http://test.com",
					1,
				).Return(errRepo)
				return mockRepo
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMock(t)
			svc := NewUrlShortenService(mockRepo)
			code, err := svc.CreateUrl(ctx, tc.url, tc.expr)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, code)
				return
			}
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Len(t, code, domain.RandomLength)
			mockRepo.AssertCalled(t, "CreateUrl", mock.Anything, code, tc.url, tc.expr)
		})
	}
}

func TestUrlShortenService_GetUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		setupMock   func(t *testing.T) *repo_mocks.UrlShortenRepository
		expectedUrl string
		expectedErr error
	}{
		{
			name: "success",
			code: "aaa",
			setupMock: func(t *testing.T) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", mock.Anything, "aaa").Return("http://test.com", nil)
				return mockRepo
			},
			expectedUrl: "http://test.com",
			expectedErr: nil,
		},
		{
			name: "repo error",
			code: "aaa",
			setupMock: func(t *testing.T) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", mock.Anything, "aaa").Return("", errRepo)
				return mockRepo
			},
			expectedErr: errRepo,
			expectedUrl: "",
		},
		{
			name: "url not found",
			code: "aaa",
			setupMock: func(t *testing.T) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", mock.Anything, "aaa").Return("", redis.Nil)
				return mockRepo
			},
			expectedUrl: "",
			expectedErr: redis.Nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMock(t)
			svc := NewUrlShortenService(mockRepo)
			url, err := svc.GetUrl(ctx, tc.code)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, url)
				return
			}
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedUrl, url)
			mockRepo.AssertCalled(t, "GetUrl", mock.Anything, tc.code)
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		length int
	}{
		{name: "zero length", length: 0},
		{name: "one char", length: 1},
		{name: "normal length", length: 10},
		{name: "long string", length: 100},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, err := generateRandomString(tc.length)

			assert.NoError(t, err)
			assert.Len(t, s, tc.length)
			assert.Regexp(t, "^[a-zA-Z0-9]*$", s)
		})
	}

	t.Run("two calls give different results", func(t *testing.T) {
		t.Parallel()
		a, err := generateRandomString(domain.RandomLength)
		assert.NoError(t, err)
		b, err := generateRandomString(domain.RandomLength)
		assert.NoError(t, err)

		assert.NotEqual(t, a, b)
	})
}
