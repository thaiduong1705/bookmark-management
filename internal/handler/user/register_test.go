package user

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/handler/response"
	svc_mocks "bookmark-management/internal/service/mocks"
	utils_mocks "bookmark-management/pkg/utils/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var serviceError = errors.New("service error")

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name          string
		req           RegisterRequest
		setupMockSvc  func(t *testing.T, ctx context.Context) *svc_mocks.UserService
		setupPassUtil func(t *testing.T) *utils_mocks.PasswordUtil
		expectedCode  int
		expectedResp  any
	}{
		{
			name: "success",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Password:    "test_password123",
				Username:    "test_username",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				user := &domain.User{
					DisplayName: "test display name",
					Email:       "test@gmail.com",
					Password:    "test_password123",
					Username:    "test_username",
				}
				mockSvc.On("CreateUser", ctx, user).Return(&domain.User{
					ID:          "ewqeqw",
					DisplayName: "test display name",
					Email:       "test@gmail.com",
					Password:    "hashed_test_password123",
					Username:    "test_username",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Validate", "test_password123").Return(true)
				return mockPass
			},
			expectedCode: http.StatusCreated,
			expectedResp: RegisterResponse{
				Message:     response.NewMessage("User registered successfully!"),
				ID:          "ewqeqw",
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Username:    "test_username",
				CreatedAt:   time.Now().Format("02-01-2006 15:04:05"),
				UpdatedAt:   time.Now().Format("02-01-2006 15:04:05"),
			},
		},
		{
			name: "service error",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Password:    "test_password123",
				Username:    "test_username",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				user := &domain.User{
					DisplayName: "test display name",
					Email:       "test@gmail.com",
					Password:    "test_password123",
					Username:    "test_username",
				}
				mockSvc.On("CreateUser", ctx, user).Return(nil, serviceError)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Validate", "test_password123").Return(true)
				return mockPass
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: response.InternalServerError,
		},
		{
			name: "user already exists",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Password:    "test_password123",
				Username:    "test_username",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				user := &domain.User{
					DisplayName: "test display name",
					Email:       "test@gmail.com",
					Password:    "test_password123",
					Username:    "test_username",
				}
				mockSvc.On("CreateUser", ctx, user).Return(nil, domain.ErrAlreadyExists)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Validate", "test_password123").Return(true)
				return mockPass
			},
			expectedCode: http.StatusConflict,
			expectedResp: response.NewErrorResponse("User already exists", "email or username is already taken"),
		},
		{
			name: "invalid user data",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Password:    "test_password123",
				Username:    "test_username",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				user := &domain.User{
					DisplayName: "test display name",
					Email:       "test@gmail.com",
					Password:    "test_password123",
					Username:    "test_username",
				}
				mockSvc.On("CreateUser", ctx, user).Return(nil, domain.ErrInvalidInput)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Validate", "test_password123").Return(true)
				return mockPass
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "invalid user data"),
		},
		{
			name: "validate password error",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Password:    "test",
				Username:    "test_username",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				mockPass.On("Validate", "test").Return(false)
				return mockPass
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "invalid Password"),
		},
		{
			name: "invalid email",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test",
				Password:    "test_password123",
				Username:    "test_username",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				return mockPass
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "Email is invalid (email)"),
		},
		{
			name: "missing username",
			req: RegisterRequest{
				DisplayName: "test display name",
				Email:       "test@gmail.com",
				Password:    "test_password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *svc_mocks.UserService {
				mockSvc := svc_mocks.NewUserService(t)
				return mockSvc
			},
			setupPassUtil: func(t *testing.T) *utils_mocks.PasswordUtil {
				mockPass := utils_mocks.NewPasswordUtil(t)
				return mockPass
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "Username is invalid (required)"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body, _ := json.Marshal(tc.req)
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request, _ = http.NewRequest(http.MethodPost, "v1/users/register", bytes.NewBuffer(body))
			mockSvc := tc.setupMockSvc(t, ginCtx)
			mockPass := tc.setupPassUtil(t)
			handler := NewHandler(mockSvc, mockPass)
			handler.Register(ginCtx)
			assert.Equal(t, tc.expectedCode, rec.Code)
			expectedResp, _ := json.Marshal(tc.expectedResp)
			assert.JSONEq(t, string(expectedResp), rec.Body.String())
		})
	}
}
