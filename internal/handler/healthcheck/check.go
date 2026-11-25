package healthcheck

import (
	"bookmark-management/internal/handler/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Check
// @Summary      Health check
// @Description  Checks that the service and its dependencies are healthy.
// @Tags         health
// @Produce      application/json
// @Success      200  {object}  CheckResponse  "Service is healthy"
// @Failure      500  {object}  response.Message  "Service is unhealthy"
// @Router       /health-check [get]
func (h *handler) Check(c *gin.Context) {
	serviceName, instanceId, err := h.srv.Check(c)
	if err != nil {
		log.Error().Err(err).Msg("health check failed")
		c.JSON(http.StatusInternalServerError, response.NewMessage("unhealthy"))
		return
	}
	c.JSON(http.StatusOK, CheckResponse{
		Message:     response.NewMessage("ok"),
		ServiceName: serviceName,
		InstanceID:  instanceId,
	})
}
