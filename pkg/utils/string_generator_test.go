package utils

import (
	"bookmark-management/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringGenerator_Generate(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		length int
	}{
		{name: "zero length", length: 0},
		{name: "one char", length: 1},
		{name: "normal length", length: 10},
		{name: "long string", length: 100},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			utilStringGen := NewStringGenerator()
			s, err := utilStringGen.Generate(tc.length)

			assert.NoError(t, err)
			assert.Len(t, s, tc.length)
			assert.Regexp(t, "^[a-zA-Z0-9]*$", s)
		})
	}

	t.Run("two calls give different results", func(t *testing.T) {
		t.Parallel()
		utilStringGen := NewStringGenerator()
		a, err := utilStringGen.Generate(domain.RandomLength)
		assert.NoError(t, err)
		b, err := utilStringGen.Generate(domain.RandomLength)
		assert.NoError(t, err)

		assert.NotEqual(t, a, b)
	})
}
