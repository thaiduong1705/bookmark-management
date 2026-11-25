package integration_tests

import (
	"bookmark-management/internal/api"
	"bookmark-management/internal/domain"
	"bookmark-management/internal/handler/user"
	redisPkg "bookmark-management/pkg/redis"
	"bookmark-management/pkg/sqldb"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		setupRequest func() (*httptest.ResponseRecorder, *http.Request)
		expectedCode int
	}{
		{
			name: "success",
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(user.RegisterRequest{
					DisplayName: "test display name",
					Email:       "test@gmail.com",
					Password:    "tEst@password123",
					Username:    "test_username",
				})
				req, _ := http.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(body))
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedCode: http.StatusCreated,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			mockDb := sqldb.InitMockDatabase(t)
			mockDb.AutoMigrate(&domain.User{})
			testApiEngine := api.NewEngine(&api.Config{}, mockRedis, mockDb)
			rec, req := tc.setupRequest()
			testApiEngine.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			var resp user.RegisterResponse
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		})
	}
}
