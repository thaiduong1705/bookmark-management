package response

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testInput mirrors a real request DTO: exported fields with binding tags.
type testInput struct {
	Url string `binding:"required,url"`
	Exp int    `binding:"gte=0"`
}

func TestParseBindingError(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		makeErr  func(t *testing.T) error
		expected ErrorResponse
	}{
		{
			name: "non-validator error only returns message",
			makeErr: func(t *testing.T) error {
				err := json.Unmarshal([]byte(`{"url":`), &struct{}{})
				require.Error(t, err)
				return err
			},
			expected: NewErrorResponse("Invalid request", ""),
		},
		{
			name: "multiple invalid fields are aggregated",
			makeErr: func(t *testing.T) error {
				return validationErr(t, &testInput{Url: "aasda", Exp: -1})
			},
			expected: NewErrorResponse("Invalid request", "Url is invalid (url) Exp is invalid (gte)"),
		},
		{
			name: "wrapped validation error is still detected",
			makeErr: func(t *testing.T) error {
				return fmt.Errorf("bind failed: %w", validationErr(t, &testInput{Url: "aasda", Exp: 1}))
			},
			expected: NewErrorResponse("Invalid request", "Url is invalid (url)"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := ParseBindingError(tc.makeErr(t))
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
