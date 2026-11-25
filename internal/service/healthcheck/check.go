package healthcheck

import (
	"context"
	"fmt"
)

// Check returns the service name and instance id, or an error if a dependency is unreachable.
func (h *service) Check(ctx context.Context) (string, string, error) {
	if err := h.repo.Check(ctx); err != nil {
		return "", "", fmt.Errorf("healthcheck: %w", err)
	}

	return h.serviceName, h.instanceId, nil
}
