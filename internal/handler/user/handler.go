package user

import (
	service "bookmark-management/internal/service/user"
	"bookmark-management/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Handler serves the user endpoints.
type Handler interface {
	Register(c *gin.Context)
}

type handler struct {
	svc service.Service
	vp  utils.PasswordUtil
}

// NewHandler creates a user Handler; vp validates the password strength.
func NewHandler(svc service.Service, vp utils.PasswordUtil) Handler {
	return &handler{svc: svc, vp: vp}
}
