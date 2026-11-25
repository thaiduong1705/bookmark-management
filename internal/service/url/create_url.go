package url

import (
	"bookmark-management/internal/domain"
	"context"
	"errors"
	"fmt"
)

// CreateUrl generates a random code and stores the URL under it.
func (svc *service) CreateUrl(ctx context.Context, url string, expr int) (string, error) {
	code, err := svc.strGenUtil.Generate(domain.RandomLength)
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}

	// check if the code exists, recursive it
	_, err = svc.repo.GetUrl(ctx, code)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", fmt.Errorf("check code exists: %w", err)
	} else if err == nil {
		return svc.CreateUrl(ctx, url, expr)
	}

	if err = svc.repo.CreateUrl(ctx, code, url, expr); err != nil {
		return "", fmt.Errorf("store url: %w", err)
	}
	return code, nil
}
