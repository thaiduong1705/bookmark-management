package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordUtil_Validate(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		password    string
		expectedRes bool
	}{
		{
			name:        "shorter than 8",
			password:    "Acx12@#",
			expectedRes: false,
		},
		{
			name:        "no uppercase letter",
			password:    "abcxzkze123@!",
			expectedRes: false,
		},
		{
			name:        "no lowercase letter",
			password:    "FHDJKASFKJFA123$#",
			expectedRes: false,
		},
		{
			name:        "no number",
			password:    "fjdajlxazEJLKAWJL@!#",
			expectedRes: false,
		},
		{
			name:        "no special",
			password:    "AJSAfdsf121",
			expectedRes: false,
		},
		{
			name:        "correct format",
			password:    "AJSDK@#!123ak",
			expectedRes: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			util := NewPasswordUtil()
			res := util.Validate(tc.password)
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestPasswordUtil_Hash(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "success",
			password:    "AJSDK@#!123ak",
			expectedErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			util := NewPasswordUtil()
			hash1, err1 := util.Hash(tc.password)
			assert.NotEqual(t, "", hash1)
			assert.NotEqual(t, tc.password, hash1)
			assert.NoError(t, err1)
			hash2, err2 := util.Hash(tc.password)
			assert.NotEqual(t, "", hash2)
			assert.NotEqual(t, tc.password, hash2)
			assert.NoError(t, err2)

			assert.NotEqual(t, hash1, hash2)
		})
	}
}

func TestPasswordUtil_Verify(t *testing.T) {
	t.Parallel()
	util := NewPasswordUtil()
	raw := "MySecret@123"
	hash, err := util.Hash(raw)
	assert.NoError(t, err)
	testcases := []struct {
		name        string
		hash        string
		password    string
		expectedRes bool
	}{
		{
			name:        "matches password",
			hash:        hash,
			password:    raw,
			expectedRes: true,
		},
		{
			name:        "empty password",
			hash:        hash,
			password:    "",
			expectedRes: false,
		},
		{
			name:        "wrong password",
			hash:        hash,
			password:    "FHDJKASFKJFA123$#",
			expectedRes: false,
		},
		{
			name:        "malformed hash return false",
			hash:        "not-a-valid-bcrypt-hash",
			password:    raw,
			expectedRes: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := util.Verify(tc.password, tc.hash)
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}
