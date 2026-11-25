package url

import (
	"context"
	"fmt"
)

// GetUrl returns the original URL for the given code.
func (svc *service) GetUrl(ctx context.Context, code string) (string, error) {
	url, err := svc.repo.GetUrl(ctx, code)
	if err != nil {
		return "", fmt.Errorf("get url: %w", err)
	}
	return url, nil
}
