package url

import (
	"bookmark-management/internal/handler/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CreateUrl
// @Summary      Create short URL
// @Description  Generates a short code for the given URL. expr is the lifetime in minutes (0 = never expires).
// @Tags         links
// @Produce      application/json
// @Param        request  body      CreateRequest  true  "URL to shorten"
// @Success      201      {object}  CreateResponse  "Generated short code"
// @Failure      400      {object}  response.ErrorResponse "Invalid request"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /v1/links/shorten [post]
func (h *handler) CreateUrl(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Str("url", req.Url).Int("exp", req.Exp).Msg("bind request failed")
		c.JSON(http.StatusBadRequest, response.ParseBindingError(err))
		return
	}

	code, err := h.svc.CreateUrl(c, req.Url, req.Exp)
	if err != nil {
		log.Error().Err(err).Msg("create url failed")
		c.JSON(http.StatusInternalServerError, response.InternalServerError)
		return
	}

	c.JSON(http.StatusCreated, CreateResponse{
		Message: response.NewMessage("Shorten URL generated successfully!"),
		Code:    code,
	})
	return
}
