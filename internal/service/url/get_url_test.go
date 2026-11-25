package url

import (
	"bookmark-management/internal/domain"
	repo_mocks "bookmark-management/internal/repository/mocks"
	utils_mocks "bookmark-management/pkg/utils/mocks"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUrlShortenService_GetUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		setupMock   func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository
		expectedUrl string
		expectedErr error
	}{
		{
			name: "success",
			code: "suc",
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", ctx, "suc").Return("http://test.com", nil)
				return mockRepo
			},
			expectedUrl: "http://test.com",
			expectedErr: nil,
		},
		{
			name: "repo error",
			code: "err",
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", ctx, "err").Return("", errRepo)
				return mockRepo
			},
			expectedErr: errRepo,
			expectedUrl: "",
		},
		{
			name: "url not found",
			code: "nofound",
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", ctx, "nofound").Return("", domain.ErrNotFound)
				return mockRepo
			},
			expectedUrl: "",
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			mockRepo := tc.setupMock(t, ctx)
			mockUtil := utils_mocks.NewStringGenerator(t)
			svc := NewService(mockRepo, mockUtil)
			url, err := svc.GetUrl(ctx, tc.code)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, url)
				return
			}
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedUrl, url)
			mockRepo.AssertCalled(t, "GetUrl", ctx, tc.code)
		})
	}
}
