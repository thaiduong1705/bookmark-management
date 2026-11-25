package url

import (
	service "bookmark-management/internal/service/url"

	"github.com/gin-gonic/gin"
)

// Handler serves the URL shortening endpoints.
type Handler interface {
	RedirectUrl(c *gin.Context)
	CreateUrl(c *gin.Context)
}

type handler struct {
	svc service.Service
}

// NewHandler creates a Handler.
func NewHandler(svc service.Service) Handler {
	return &handler{
		svc: svc,
	}
}
