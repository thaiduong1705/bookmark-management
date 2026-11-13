package service

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/repository"
	"crypto/rand"
	"math/big"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --name=UrlShortenService --filename urlshorten.go --outpkg svc_mocks

// UrlShortenService contains the URL shortening business logic.
type UrlShortenService interface {
	GetUrl(c *gin.Context, code string) (string, error)
	CreateUrl(c *gin.Context, url string, expr int) (string, error)
}

type urlShortenService struct {
	repo repository.UrlShortenRepository
}

// NewUrlShortenService creates a UrlShortenService.
func NewUrlShortenService(repo repository.UrlShortenRepository) UrlShortenService {
	return &urlShortenService{
		repo: repo,
	}
}

// GetUrl returns the original URL for the given code.
func (svc *urlShortenService) GetUrl(c *gin.Context, code string) (string, error) {
	return svc.repo.GetUrl(c, code)
}

// CreateUrl generates a random code and stores the URL under it.
func (svc *urlShortenService) CreateUrl(c *gin.Context, url string, expr int) (string, error) {
	code, err := generateRandomString(domain.RandomLength)
	if err != nil {
		return "", err
	}

	if err = svc.repo.CreateUrl(c, code, url, expr); err != nil {
		return "", err
	}
	return code, nil
}

// generateRandomString returns a cryptographically random alphanumeric string.
func generateRandomString(length int) (string, error) {
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
