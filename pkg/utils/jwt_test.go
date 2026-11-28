package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func writeTestKeys(t *testing.T) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	assert.NoError(t, err)

	dir := t.TempDir()
	privatePath := filepath.Join(dir, "private.pem")
	publicPath := filepath.Join(dir, "public.pem")

	assert.NoError(t, os.WriteFile(privatePath, pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
	}), 0o600))
	assert.NoError(t, os.WriteFile(publicPath, pem.EncodeToMemory(&pem.Block{
		Type: "PUBLIC KEY", Bytes: pubBytes,
	}), 0o600))

	return privatePath, publicPath
}

func TestJwtUtil_Generate(t *testing.T) {
	testcases := []struct {
		name        string
		sub         jwt.MapClaims
		expectedErr error
	}{
		{
			name: "success",
			sub: jwt.MapClaims{
				"sub": "1",
			},
		},
		{
			name: "fail",
			sub: jwt.MapClaims{
				"x": make(chan int),
			},
			expectedErr: errors.New("json: unsupported type: chan int"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			jwtUtil, err := NewJwtUtil(writeTestKeys(t))
			assert.NoError(t, err)
			token, err := jwtUtil.Generate(tc.sub)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Empty(t, token)
				return
			}
			assert.NoError(t, err)
			splits := strings.Split(token, ".")
			assert.Len(t, splits, 3)
		})
	}
}

func TestJwtUtil_Verify(t *testing.T) {
	t.Parallel()

	newJwtUtil := func(t *testing.T) JwtUtil {
		t.Helper()
		jwtUtil, err := NewJwtUtil(writeTestKeys(t))
		assert.NoError(t, err)
		return jwtUtil
	}

	testcases := []struct {
		name           string
		setup          func(t *testing.T) (JwtUtil, string)
		expectedClaims jwt.MapClaims
		expectedErr    error
	}{
		{
			name: "success",
			setup: func(t *testing.T) (JwtUtil, string) {
				jwtUtil := newJwtUtil(t)
				token, err := jwtUtil.Generate(jwt.MapClaims{"sub": "1"})
				assert.NoError(t, err)
				return jwtUtil, token
			},
			expectedClaims: jwt.MapClaims{"sub": "1"},
		},
		{
			name: "expired token",
			setup: func(t *testing.T) (JwtUtil, string) {
				jwtUtil := newJwtUtil(t)
				token, err := jwtUtil.Generate(jwt.MapClaims{"sub": "1", "exp": 100})
				assert.NoError(t, err)
				return jwtUtil, token
			},
			expectedErr: jwt.ErrTokenExpired,
		},
		{
			name: "trash string",
			setup: func(t *testing.T) (JwtUtil, string) {
				return newJwtUtil(t), "abc"
			},
			expectedErr: jwt.ErrTokenMalformed,
		},
		{
			name: "wrong algo",
			setup: func(t *testing.T) (JwtUtil, string) {
				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "1"}).
					SignedString([]byte("secret"))
				assert.NoError(t, err)
				return newJwtUtil(t), token
			},
			expectedErr: jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			jwtUtil, token := tc.setup(t)

			claims, err := jwtUtil.Validate(token)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, claims)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedClaims, claims)
		})
	}
}
