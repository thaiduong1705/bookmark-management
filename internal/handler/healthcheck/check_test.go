package healthcheck

import (
	"bookmark-management/internal/handler/response"
	healthcheck_svc_mocks "bookmark-management/internal/service/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		svcName        string
		svcInstanceId  string
		svcErr         error
		expectedStatus int
		expectedBody   any
	}{
		{
			name:           "OK",
			svcName:        "bookmark_service",
			svcInstanceId:  "instance_1",
			expectedStatus: http.StatusOK,
			expectedBody: CheckResponse{
				Message:     response.NewMessage("ok"),
				ServiceName: "bookmark_service",
				InstanceID:  "instance_1",
			},
		},
		{
			name:           "redis down",
			svcErr:         errors.New("redis down"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   response.NewMessage("unhealthy"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = httptest.NewRequest("GET", "/health-check", nil)

			mockSvc := healthcheck_svc_mocks.NewHealthCheck(t)
			mockSvc.On("Check", ginCtx).Return(tc.svcName, tc.svcInstanceId, tc.svcErr)
			NewHandler(mockSvc).Check(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			expectedBody, _ := json.Marshal(tc.expectedBody)
			assert.JSONEq(t, string(expectedBody), rec.Body.String())
		})
	}
}
