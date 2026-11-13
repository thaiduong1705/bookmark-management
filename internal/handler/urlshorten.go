package handler

import (
	"bookmark-management/internal/domain"
	"bookmark-management/internal/service"
	"errors"
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
	Url  string `json:"url" binding:"required,url"`
	Expr int    `json:"expr" binding:"gte=0"`
}

// UrlShortenQuery is the path parameter of a redirect request.
type UrlShortenQuery struct {
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
// @Failure      400   {object}  map[string]string  "Invalid query"
// @Failure      404   {object}  map[string]string  "URL not found"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Router       /v1/links/shorten/{code} [get]
func (h *urlShortenHandler) RedirectUrl(c *gin.Context) {
	var query UrlShortenQuery
	if err := c.ShouldBindUri(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	url, err := h.svc.GetUrl(c, query.Code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
			return
		}
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
// @Failure      400      {object}  map[string]string  "Invalid request"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /v1/links/shorten [post]
func (h *urlShortenHandler) CreateUrl(c *gin.Context) {
	var req UrlShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	code, err := h.svc.CreateUrl(c, req.Url, req.Expr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": code})
	return
}
