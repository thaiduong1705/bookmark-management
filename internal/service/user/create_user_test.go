package user

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

var hashError = errors.New("hash error")
var createUserError = errors.New("create user error")

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name          string
		inputModel    *domain.User
		expectedErr   error
		setupMockRepo func(t *testing.T, ctx context.Context) *repo_mocks.UserRepository
		setupMockPass func(t *testing.T, ctx context.Context) *utils_mocks.PasswordUtil
	}{
		{
			name: "success",
			inputModel: &domain.User{
				Username:    "test",
				Password:    "test",
				Email:       "test@gmail.com",
				DisplayName: "test display name",
			},
			expectedErr: nil,
			setupMockRepo: func(t *testing.T, ctx context.Context) *repo_mocks.UserRepository {
				mockRepo := repo_mocks.NewUserRepository(t)
				user := &domain.User{
					Username:    "test",
					Password:    "hashed_test",
					Email:       "test@gmail.com",
					DisplayName: "test display name",
				}
				mockRepo.On("CreateUser", ctx, user).Return(&domain.User{
					Username:    "test",
					Password:    "hashed_test",
					Email:       "test@gmail.com",
					DisplayName: "test display name",
				}, nil)
				return mockRepo
			},
			setupMockPass: func(t *testing.T, ctx context.Context) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Hash", "test").Return("hashed_test", nil)
				return mockPass
			},
		},
		{
			name: "hashed got error",
			inputModel: &domain.User{
				Username:    "test",
				Password:    "test",
				Email:       "test@gmail.com",
				DisplayName: "test display name",
			},
			setupMockRepo: func(t *testing.T, ctx context.Context) *repo_mocks.UserRepository {
				mockRepo := repo_mocks.NewUserRepository(t)
				return mockRepo
			},
			setupMockPass: func(t *testing.T, ctx context.Context) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Hash", "test").Return("", hashError)
				return mockPass
			},
			expectedErr: hashError,
		},
		{
			name: "create got error",
			inputModel: &domain.User{
				Username:    "test",
				Password:    "test",
				Email:       "test@gmail.com",
				DisplayName: "test display name",
			},
			setupMockRepo: func(t *testing.T, ctx context.Context) *repo_mocks.UserRepository {
				mockRepo := repo_mocks.NewUserRepository(t)
				user := &domain.User{
					Username:    "test",
					Password:    "hashed_test",
					Email:       "test@gmail.com",
					DisplayName: "test display name",
				}
				mockRepo.On("CreateUser", ctx, user).Return(nil, createUserError)
				return mockRepo
			},
			setupMockPass: func(t *testing.T, ctx context.Context) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Hash", "test").Return("hashed_test", nil)
				return mockPass
			},
			expectedErr: createUserError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			mockRepo := tc.setupMockRepo(t, ctx)
			mockPass := tc.setupMockPass(t, ctx)
			svc := NewService(mockRepo, mockPass)
			user, err := svc.CreateUser(ctx, tc.inputModel)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, user)
		})
	}
}
