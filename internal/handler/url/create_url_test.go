package url

import (
	"bookmark-management/internal/handler/response"
	svc_mocks "bookmark-management/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var svcError = errors.New("svc error")

func TestUrlShortenHandler_CreateUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		req          CreateRequest
		setupMock    func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService
		expectedCode int
		expectedResp any
	}{
		{
			name: "success",
			req: CreateRequest{
				Url: "http://test.com",
				Exp: 1,
			},
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("CreateUrl", ctx, "http://test.com", 1).Return("testCode", nil)
				return mockSvc
			},
			expectedCode: http.StatusCreated,
			expectedResp: CreateResponse{
				Message: response.NewMessage("Shorten URL generated successfully!"),
				Code:    "testCode",
			},
		},
		{
			name: "invalid url",
			req: CreateRequest{
				Url: "aasda",
				Exp: 1,
			},
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				return mockSvc
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "Url is invalid (url)"),
		},
		{
			name: "invalid expr",
			req: CreateRequest{
				Url: "http://test.com",
				Exp: -1,
			},
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				return mockSvc
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "Exp is invalid (gte)"),
		},
		{
			name: "svc error",
			req: CreateRequest{
				Url: "http://test.com",
				Exp: 1,
			},
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("CreateUrl", ctx, "http://test.com", 1).Return("", svcError)
				return mockSvc
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: response.InternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body, _ := json.Marshal(tc.req)
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/link/shorten", bytes.NewReader(body))
			mockSvc := tc.setupMock(t, ginCtx)
			handler := NewHandler(mockSvc)
			handler.CreateUrl(ginCtx)
			assert.Equal(t, tc.expectedCode, rec.Code)
			expectedResp, _ := json.Marshal(tc.expectedResp)
			assert.JSONEq(t, string(expectedResp), rec.Body.String())
		})
	}
}
