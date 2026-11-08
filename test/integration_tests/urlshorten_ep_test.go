package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"shorten-url/internal/api"
	"shorten-url/internal/domain"
	"shorten-url/internal/handler"
	redisPkg "shorten-url/pkg/redis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlShorten_CreateUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		setupRequest func() (*httptest.ResponseRecorder, *http.Request)
		expectedCode int
		expectedResp string
	}{
		{
			name: "success",
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				body, _ := json.Marshal(handler.UrlShortenRequest{
					Url:  "http://test.com",
					Expr: 1,
				})
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			defer mockRedis.Close()
			testApiEngine := api.NewEngine(&api.Config{}, mockRedis)
			rec, req := tc.setupRequest()
			testApiEngine.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedResp != "" {
				assert.JSONEq(t, tc.expectedResp, rec.Body.String())
				return
			}
			var resp struct {
				Code string `json:"code"`
			}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Len(t, resp.Code, domain.RandomLength)
		})
	}
}

func TestUrlShorten_GetUrl(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name         string
		code         string
		url          string
		setupRequest func() (*httptest.ResponseRecorder, *http.Request)
		expectedCode int
		expectedResp string
	}{
		{
			name: "success",
			code: "test",
			url:  "http://test.com",
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				target := "/v1/links/shorten?code=" + url.QueryEscape("test")
				req := httptest.NewRequest(http.MethodGet, target, nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedCode: http.StatusOK,
			expectedResp: `{"url":"http://test.com"}`,
		},
		{
			name: "code not found",
			code: "test",
			url:  "http://test.com",
			setupRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				target := "/v1/links/shorten?code=" + url.QueryEscape("aaa")
				req := httptest.NewRequest(http.MethodGet, target, nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"error": "url not found"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRedis := redisPkg.InitMockRedis(t)
			mockRedis.Set(t.Context(), tc.code, tc.url, 0)
			defer mockRedis.Close()
			testApiEngine := api.NewEngine(&api.Config{}, mockRedis)
			rec, req := tc.setupRequest()
			testApiEngine.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedResp, rec.Body.String())
		})
	}
}
