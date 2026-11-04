package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		serviceName string
		instanceID  string
	}{
		{"normal case", "normally", "normally"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotName, gotInstanceID := NewHealthCheck(tc.serviceName, tc.instanceID).Check()
			assert.Equal(t, tc.serviceName, gotName)
			assert.Equal(t, tc.instanceID, gotInstanceID)
		})
	}
}
