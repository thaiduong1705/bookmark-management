package handler

import (
	"net/http"
	"shorten-url/internal/service"

	"github.com/gin-gonic/gin"
)

type UrlShortenHandler interface {
	GetUrl(c *gin.Context)
	CreateUrl(c *gin.Context)
}

type urlShortenHandler struct {
	svc service.UrlShortenService
}

func NewUrlShortenHandler(svc service.UrlShortenService) UrlShortenHandler {
	return &urlShortenHandler{
		svc: svc,
	}
}

type urlShortenRequest struct {
	Url  string `json:"url"`
	Expr int    `json:"expr"`
}

func (h *urlShortenHandler) GetUrl(c *gin.Context) {
	return
}

func (h *urlShortenHandler) CreateUrl(c *gin.Context) {
	var req urlShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	code, err := h.svc.CreateUrl(c, req.Url, req.Expr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
	return
}
