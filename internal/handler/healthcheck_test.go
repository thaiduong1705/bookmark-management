package handler

import (
	healthcheck_svc_mocks "bookmark-management/internal/service/mocks"
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
		svcErr         error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "OK",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"ok","redis":"ok"}`,
		},
		{
			name:           "redis down",
			svcErr:         errors.New("redis down"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"unhealthy","redis":"down"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = httptest.NewRequest("GET", "/health-check", nil)

			mockSvc := healthcheck_svc_mocks.NewHealthCheck(t)
			mockSvc.On("Check", ginCtx).Return(tc.svcErr)
			NewHealthCheck(mockSvc).Check(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
