package url

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/handler/response"
	svc_mocks "bookmark-management/internal/service/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUrlShortenHandler_RedirectUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		code         string
		setupMock    func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService
		expectedCode int
		expectedLoc  string
		expectedResp any
	}{
		{
			name: "success",
			code: "test",
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("GetUrl", ctx, "test").Return("http://test.com", nil)
				return mockSvc
			},
			expectedCode: http.StatusMovedPermanently,
			expectedLoc:  "http://test.com",
		},
		{
			name: "invalid code",
			code: "",
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				return svc_mocks.NewUrlShortenService(t)
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: response.NewErrorResponse("Invalid request", "Code is invalid (required)"),
		},
		{
			name: "code not found",
			code: "aaa",
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("GetUrl", ctx, "aaa").Return("", domain.ErrNotFound)
				return mockSvc
			},
			expectedCode: http.StatusNotFound,
			expectedResp: response.NewErrorResponse("URL not found", ""),
		},
		{
			name: "service error",
			code: "test",
			setupMock: func(t *testing.T, ctx *gin.Context) *svc_mocks.UrlShortenService {
				mockSvc := svc_mocks.NewUrlShortenService(t)
				mockSvc.On("GetUrl", ctx, "test").Return("", svcError)
				return mockSvc
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: response.InternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			ginCtx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/shorten/"+url.PathEscape(tc.code), nil)
			ginCtx.Params = gin.Params{{Key: "code", Value: tc.code}}
			mockSvc := tc.setupMock(t, ginCtx)
			handler := NewHandler(mockSvc)
			handler.RedirectUrl(ginCtx)
			ginCtx.Writer.WriteHeaderNow()
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedLoc, rec.Header().Get("Location"))
			if tc.expectedResp != nil {
				expectedResp, _ := json.Marshal(tc.expectedResp)
				assert.JSONEq(t, string(expectedResp), rec.Body.String())
			}
		})
	}
}
