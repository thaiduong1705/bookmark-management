package user

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/handler/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Register
// @Summary      Register account
// @Description  User registers.
// @Tags         users
// @Produce      application/json
// @Param        request  body      RegisterRequest  true  "Register payload"
// @Success      201      {object}  RegisterResponse  "Generated user"
// @Failure      400      {object}  response.ErrorResponse "Invalid request"
// @Failure      409      {object}  response.ErrorResponse  "User already exists"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /v1/users/register [post]
func (h *handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("bind request failed")
		c.JSON(http.StatusBadRequest, response.ParseBindingError(err))
		return
	}

	if ok := h.vp.Validate(req.Password); !ok {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request", "invalid Password"))
		return
	}

	var user *domain.User
	user, err := h.svc.CreateUser(c, &domain.User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Username:    req.Username,
		Password:    req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAlreadyExists):
			log.Warn().Err(err).Str("username", req.Username).Msg("user already exists")
			c.JSON(http.StatusConflict, response.NewErrorResponse("User already exists", "email or username is already taken"))
		case errors.Is(err, domain.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request", "invalid user data"))
		default:
			log.Error().Err(err).Msg("create user failed")
			c.JSON(http.StatusInternalServerError, response.InternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Message:     response.NewMessage("User registered successfully!"),
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Username:    user.Username,
		CreatedAt:   user.CreatedAt.Format("02-01-2006 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("02-01-2006 15:04:05"),
	})
}
