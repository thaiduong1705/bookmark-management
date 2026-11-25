package url

import (
	"bookmark-management/internal/domain"
	repo_mocks "bookmark-management/internal/repository/mocks"
	utils_mocks "bookmark-management/pkg/utils/mocks"
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var errRepo = errors.New("repo error")
var errCode = errors.New("code gen error")

func TestUrlShortenService_CreateUrl(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name            string
		url             string
		expr            int
		setupMock       func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository
		setupMockStrGen func(t *testing.T, length int) *utils_mocks.StringGenerator
		expectedErr     error
		expectedCode    string
	}{
		{
			name: "success",
			url:  "http://test.com",
			expr: 1,
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", ctx, "abc1234").Return("", domain.ErrNotFound)
				mockRepo.On("CreateUrl",
					ctx,
					"abc1234",
					"http://test.com",
					1,
				).Return(nil)
				return mockRepo
			},
			setupMockStrGen: func(t *testing.T, length int) *utils_mocks.StringGenerator {
				mockUtil := utils_mocks.NewStringGenerator(t)
				mockUtil.On("Generate", length).Return("abc1234", nil)
				return mockUtil
			},
			expectedErr:  nil,
			expectedCode: "abc1234",
		},
		{
			name: "code collision then retry",
			url:  "http://test.com",
			expr: 1,
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", ctx, "abc1234").Return("http://other.com", nil).Once()
				mockRepo.On("GetUrl", ctx, "xyz5678").Return("", domain.ErrNotFound).Once()
				mockRepo.On("CreateUrl", ctx, "xyz5678", "http://test.com", 1).Return(nil).Once()
				return mockRepo
			},
			setupMockStrGen: func(t *testing.T, length int) *utils_mocks.StringGenerator {
				mockUtil := utils_mocks.NewStringGenerator(t)
				mockUtil.On("Generate", length).Return("abc1234", nil).Once()
				mockUtil.On("Generate", length).Return("xyz5678", nil).Once()
				return mockUtil
			},
			expectedErr:  nil,
			expectedCode: "xyz5678",
		},
		{
			name: "code gen error",
			url:  "http://test.com",
			expr: 1,
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				return mockRepo
			},
			setupMockStrGen: func(t *testing.T, length int) *utils_mocks.StringGenerator {
				mockUtil := utils_mocks.NewStringGenerator(t)
				mockUtil.On("Generate", length).Return("", errCode)
				return mockUtil
			},
			expectedErr:  errCode,
			expectedCode: "",
		},
		{
			name: "repo error",
			url:  "http://test.com",
			expr: 1,
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.UrlShortenRepository {
				mockRepo := repo_mocks.NewUrlShortenRepository(t)
				mockRepo.On("GetUrl", ctx, "repoerr").Return("", domain.ErrNotFound)
				mockRepo.On("CreateUrl",
					ctx,
					"repoerr",
					"http://test.com",
					1,
				).Return(errRepo)
				return mockRepo
			},
			setupMockStrGen: func(t *testing.T, length int) *utils_mocks.StringGenerator {
				mockUtil := utils_mocks.NewStringGenerator(t)
				mockUtil.On("Generate", length).Return("repoerr", nil)
				return mockUtil
			},
			expectedErr:  errRepo,
			expectedCode: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			mockRepo := tc.setupMock(t, ctx)
			mockStrGen := tc.setupMockStrGen(t, domain.RandomLength)
			svc := NewService(mockRepo, mockStrGen)
			code, err := svc.CreateUrl(ctx, tc.url, tc.expr)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, code)
				return
			}
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedCode, code)
			mockRepo.AssertCalled(t, "CreateUrl", ctx, code, tc.url, tc.expr)
		})
	}
}
