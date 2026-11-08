package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	healthcheck_svc_mocks "shorten-url/internal/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			mockSvc.On("Check", mock.Anything).Return(tc.svcErr)
			NewHealthCheck(mockSvc).Check(ginCtx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
