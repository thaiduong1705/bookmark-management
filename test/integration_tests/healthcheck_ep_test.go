package integration_tests

import (
	"net/http"
	"net/http/httptest"
	"shorten-url/internal/api"
	"shorten-url/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		cfg            config.Config
		setupRequest   func() (*httptest.ResponseRecorder, *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "OK",
			cfg: config.Config{
				ServiceName: "svc_1",
				InstanceID:  "ins_1",
			},
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health-check", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"instance_id":"ins_1","message":"ok","service_name":"svc_1"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testApiEngine := api.NewEngine(&tc.cfg)
			rec, req := tc.setupRequest()
			testApiEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
