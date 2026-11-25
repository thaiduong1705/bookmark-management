package healthcheck

import "bookmark-management/internal/handler/response"

// CheckResponse is the body returned when the service is healthy.
type CheckResponse struct {
	response.Message
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}
