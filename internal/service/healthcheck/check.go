package healthcheck

import (
	"context"
	"fmt"
)

// Check pings Redis and returns an error if it is unreachable.
func (h *service) Check(ctx context.Context) (string, string, error) {
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		return "", "", fmt.Errorf("redis ping: %w", err)
	}
	return h.serviceName, h.instanceId, nil
}
