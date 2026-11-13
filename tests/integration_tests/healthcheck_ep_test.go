package integration_tests

import (
	"bookmark-management/internal/api"
	redisPkg "bookmark-management/pkg/redis"
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
		expectedBody   string
	}{
		{
			name: "OK",
			cfg:  api.Config{},
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health-check", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"ok","redis":"ok"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			defer mockRedis.Close()
			testApiEngine := api.NewEngine(&tc.cfg, mockRedis)
			rec, req := tc.setupRequest()
			testApiEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
