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
	GetUrl(c *gin.Context)
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

// UrlShortenQuery is the query of a get-URL request.
type UrlShortenQuery struct {
	Code string `form:"code" binding:"required"`
}

// GetUrl
// @Summary      Get original URL
// @Description  Returns the original URL stored under the given short code.
// @Tags         links
// @Produce      application/json
// @Param        code  query     string             true  "Short code"
// @Success      200   {object}  map[string]string  "Original URL"
// @Failure      400   {object}  map[string]string  "Invalid query"
// @Failure      404   {object}  map[string]string  "URL not found"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Router       /v1/links/shorten [get]
func (h *urlShortenHandler) GetUrl(c *gin.Context) {
	var query UrlShortenQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	url, err := h.svc.GetUrl(c.Request.Context(), query.Code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
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

	code, err := h.svc.CreateUrl(c.Request.Context(), req.Url, req.Expr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": code})
	return
}
