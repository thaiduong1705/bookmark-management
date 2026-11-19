package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBindingError(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		makeErr  func(t *testing.T) error
		expected FieldErrorResponse
	}{
		{
			name: "non-validator error only returns message",
			makeErr: func(t *testing.T) error {
				err := json.Unmarshal([]byte(`{"url":`), &UrlShortenRequest{})
				require.Error(t, err)
				return err
			},
			expected: FieldErrorResponse{Message: "Invalid request"},
		},
		{
			name: "multiple invalid fields are aggregated",
			makeErr: func(t *testing.T) error {
				return validationErr(t, &UrlShortenRequest{Url: "aasda", Exp: -1})
			},
			expected: FieldErrorResponse{Message: "Invalid request", Detail: "Url is invalid (url) Exp is invalid (gte)"},
		},
		{
			name: "wrapped validation error is still detected",
			makeErr: func(t *testing.T) error {
				return fmt.Errorf("bind failed: %w", validationErr(t, &UrlShortenRequest{Url: "aasda", Exp: 1}))
			},
			expected: FieldErrorResponse{Message: "Invalid request", Detail: "Url is invalid (url)"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := parseBindingError(tc.makeErr(t))
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// validationErr validates obj with gin's validator.
func validationErr(t *testing.T, obj any) error {
	t.Helper()
	err := binding.Validator.ValidateStruct(obj)
	require.Error(t, err)
	return err
}
