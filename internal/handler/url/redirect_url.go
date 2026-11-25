package url

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/handler/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RedirectUrl
// @Summary      Redirect to original URL
// @Description  Looks up the short code and redirects (301) to the original URL via the Location header.
// @Tags         links
// @Produce      application/json
// @Param        code  path      string             true  "Short code"
// @Success      301   {string}  string             "Redirect to original URL (see Location header)"
// @Header       301   {string}  Location           "Original URL"
// @Failure      400   {object}  response.ErrorResponse "Invalid query"
// @Failure      404   {object}  response.ErrorResponse  "URL not found"
// @Failure      500   {object}  response.ErrorResponse  "Internal server error"
// @Router       /v1/links/shorten/{code} [get]
func (h *handler) RedirectUrl(c *gin.Context) {
	var param CodeParam
	if err := c.ShouldBindUri(&param); err != nil {
		log.Warn().Err(err).Str("code", param.Code).Msg("bind request failed")
		c.JSON(http.StatusBadRequest, response.ParseBindingError(err))
		return
	}

	url, err := h.svc.GetUrl(c, param.Code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse("URL not found", ""))
			return
		}
		log.Error().Err(err).Str("code", param.Code).Msg("get url failed")
		c.JSON(http.StatusInternalServerError, response.InternalServerError)
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)
}
