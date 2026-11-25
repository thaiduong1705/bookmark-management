package utils

import (
	"crypto/rand"
	"math/big"
)

// StringGenerator generates random strings.
type StringGenerator interface {
	// Generate returns a random alphanumeric string of the given length.
	Generate(length int) (string, error)
}

type stringGenerator struct{}

// NewStringGenerator creates a StringGenerator that uses a cryptographically secure random source.
func NewStringGenerator() StringGenerator {
	return &stringGenerator{}
}

// Generate returns a cryptographically random alphanumeric string.
func (g *stringGenerator) Generate(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
