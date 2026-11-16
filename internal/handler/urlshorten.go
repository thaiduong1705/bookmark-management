package handler

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/service"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UrlShortenHandler serves the URL shortening endpoints.
type UrlShortenHandler interface {
	RedirectUrl(c *gin.Context)
	CreateUrl(c *gin.Context)
}

type urlShortenHandler struct {
	svc service.UrlShortenService
}

// NewUrlShortenHandler creates a UrlShortenHandler.
func NewUrlShortenHandler(svc service.UrlShortenService) UrlShortenHandler {
	return &urlShortenHandler{
		svc: svc,
	}
}

// UrlShortenRequest is the body of a create-short-URL request; Expr is in minutes (0 = never expires).
type UrlShortenRequest struct {
	Url string `json:"url" binding:"required,url"`
	Exp int    `json:"exp" binding:"gte=0,lte=84600"`
}

// UrlShortenParam is the path parameter of a redirect request.
type UrlShortenParam struct {
	Code string `uri:"code" binding:"required"`
}

// RedirectUrl
// @Summary      Redirect to original URL
// @Description  Looks up the short code and redirects (301) to the original URL via the Location header.
// @Tags         links
// @Produce      application/json
// @Param        code  path      string             true  "Short code"
// @Success      301   {string}  string             "Redirect to original URL (see Location header)"
// @Header       301   {string}  Location           "Original URL"
// @Failure      400   {object}  FieldErrorResponse "Invalid query"
// @Failure      404   {object}  map[string]string  "URL not found"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Router       /v1/links/shorten/{code} [get]
func (h *urlShortenHandler) RedirectUrl(c *gin.Context) {
	var param UrlShortenParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, parseBindingError(err))
		return
	}

	url, err := h.svc.GetUrl(c, param.Code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
			return
		}
		slog.ErrorContext(c, "failed to get url", "code", param.Code, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)
}

// CreateUrl
// @Summary      Create short URL
// @Description  Generates a short code for the given URL. expr is the lifetime in minutes (0 = never expires).
// @Tags         links
// @Produce      application/json
// @Param        request  body      UrlShortenRequest  true  "URL to shorten"
// @Success      201      {object}  map[string]string  "Generated short code"
// @Failure      400      {object}  FieldErrorResponse "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /v1/links/shorten [post]
func (h *urlShortenHandler) CreateUrl(c *gin.Context) {
	var req UrlShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, parseBindingError(err))
		return
	}

	code, err := h.svc.CreateUrl(c, req.Url, req.Exp)
	if err != nil {
		slog.ErrorContext(c, "failed to create url", "exp", req.Exp, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": code, "message": "Shorten URL generated successfully!"})
	return
}
