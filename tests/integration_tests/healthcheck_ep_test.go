package integration_tests

import (
	"bookmark-management/internal/api"
	"bookmark-management/internal/handler/healthcheck"
	"bookmark-management/internal/handler/response"
	redisPkg "bookmark-management/pkg/redis"
	"bookmark-management/pkg/sqldb"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		cfg            api.Config
		setupRequest   func() (*httptest.ResponseRecorder, *http.Request)
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "OK",
			cfg:  api.Config{ServiceName: "bookmark_service", InstanceId: "instance_1"},
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health-check", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus: http.StatusOK,
			expectedBody: healthcheck.CheckResponse{
				Message:     response.NewMessage("ok"),
				ServiceName: "bookmark_service",
				InstanceID:  "instance_1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			mockDb := sqldb.InitMockDatabase(t)
			testApiEngine := api.NewEngine(&tc.cfg, mockRedis, mockDb)
			rec, req := tc.setupRequest()
			testApiEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			expectedBody, _ := json.Marshal(tc.expectedBody)
			assert.JSONEq(t, string(expectedBody), rec.Body.String())
		})
	}
}
