package handler

import (
	"net/http"
	"net/http/httptest"
	healthcheck_svc_mocks "shorten-url/internal/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serviceName    string
		instanceID     string
		setupMockSvc   func(t *testing.T, serviceName, instanceID string) *healthcheck_svc_mocks.HealthCheck
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "OK",
			serviceName: "service_test",
			instanceID:  "instance_test",
			setupMockSvc: func(t *testing.T, serviceName, instanceID string) *healthcheck_svc_mocks.HealthCheck {
				mockSvc := healthcheck_svc_mocks.NewHealthCheck(t)
				mockSvc.On("Check").Return(serviceName, instanceID)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"instance_id":"instance_test","message":"ok","service_name":"service_test"}`,
		},

		{
			name:        "OK2",
			serviceName: "svc_2_test",
			instanceID:  "ins_2_test",
			setupMockSvc: func(t *testing.T, serviceName, instanceID string) *healthcheck_svc_mocks.HealthCheck {
				mockSvc := healthcheck_svc_mocks.NewHealthCheck(t)
				mockSvc.On("Check").Return(serviceName, instanceID)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"instance_id":"ins_2_test","message":"ok","service_name":"svc_2_test"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)

			mockSvc := tc.setupMockSvc(t, tc.serviceName, tc.instanceID)
			testHandler := NewHealthCheck(mockSvc)
			testHandler.Check(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
