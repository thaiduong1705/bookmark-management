package handler

import (
	"bookmark-management/internal/domain"
	svc_mocks "bookmark-management/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var svcError = errors.New("svc error")

func TestUrlShortenHandler_CreateUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		req          UrlShortenRequest
		setupMock    func(t *testing.T) *svc_mocks.UrlShortenService
		expectedCode int
		expectedResp string
	}{
		{
			name: "success",
			req: UrlShortenRequest{
				Url:  "http://test.com",
				Expr: 1,
			},
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("CreateUrl", mock.Anything, "http://test.com", 1).Return("testCode", nil)
				return mockSvc
			},
			expectedCode: http.StatusCreated,
			expectedResp: `{"code":"testCode"}`,
		},
		{
			name: "invalid url",
			req: UrlShortenRequest{
				Url:  "aasda",
				Expr: 1,
			},
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				return mockSvc
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: `{"error": "invalid request"}`,
		},
		{
			name: "invalid expr",
			req: UrlShortenRequest{
				Url:  "http://test.com",
				Expr: -1,
			},
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				return mockSvc
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: `{"error": "invalid request"}`,
		},
		{
			name: "svc error",
			req: UrlShortenRequest{
				Url:  "http://test.com",
				Expr: 1,
			},
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("CreateUrl", mock.Anything, "http://test.com", 1).Return("", svcError)
				return mockSvc
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"error": "internal server error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			body, _ := json.Marshal(tc.req)
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = httptest.NewRequest(http.MethodPost, "/v1/link/shorten", bytes.NewReader(body))
			mockSvc := tc.setupMock(t)
			handler := NewUrlShortenHandler(mockSvc)
			handler.CreateUrl(ginCtx)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedResp, rec.Body.String())
		})
	}
}

func TestUrlShortenHandler_GetUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		code         string
		setupMock    func(t *testing.T) *svc_mocks.UrlShortenService
		expectedCode int
		expectedResp string
	}{
		{
			name: "success",
			code: "tests",
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("GetUrl", mock.Anything, "tests").Return("http://test.com", nil)
				return mockSvc
			},
			expectedCode: http.StatusOK,
			expectedResp: `{"url":"http://tests.com"}`,
		},
		{
			name: "invalid code",
			code: "",
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				return svc_mocks.NewUrlShortenService(t) // service không được gọi
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: `{"error": "invalid query"}`,
		},
		{
			name: "code not found",
			code: "aaa",
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("GetUrl", mock.Anything, "aaa").Return("", domain.ErrNotFound)
				return mockSvc
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"error": "url not found"}`,
		},
		{
			name: "service error",
			code: "tests",
			setupMock: func(t *testing.T) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("GetUrl", mock.Anything, "tests").Return("", svcError)
				return mockSvc
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"error": "internal server error"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			target := "/v1/link/shorten?code=" + url.QueryEscape(tc.code)
			ginCtx.Request = httptest.NewRequest(http.MethodGet, target, nil)
			mockSvc := tc.setupMock(t)
			handler := NewUrlShortenHandler(mockSvc)
			handler.GetUrl(ginCtx)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedResp, rec.Body.String())
		})
	}
}
