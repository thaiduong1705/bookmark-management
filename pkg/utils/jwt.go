package utils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtil interface {
	Generate(jwtContent jwt.MapClaims) (string, error)
	Validate(tokenStr string) (jwt.MapClaims, error)
}

type jwtUtil struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJwtUtil(privatePath, publicPath string) (JwtUtil, error) {
	pbKeyData, err := os.ReadFile(publicPath)
	if err != nil {
		return nil, err
	}
	prKeyData, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(prKeyData)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pbKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtUtil{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (j *jwtUtil) Generate(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	tokenStr, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *jwtUtil) Validate(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.publicKey, nil
	},
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
