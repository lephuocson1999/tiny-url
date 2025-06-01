package api

import (
	"net/http"
	"time"

	"tiny-url/internal/app"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service app.URLShortenerService
}

func NewHandler(service app.URLShortenerService) *Handler {
	return &Handler{service: service}
}

type ShortenRequest struct {
	URL       string     `json:"url" binding:"required,url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ShortenResponse struct {
	ShortURL  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type StatsResponse struct {
	OriginalURL  string     `json:"original_url"`
	ShortURL     string     `json:"short_url"`
	AccessCount  int64      `json:"access_count"`
	CreatedAt    time.Time  `json:"created_at"`
	LastAccessed *time.Time `json:"last_accessed"`
}

func (h *Handler) ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	code, err := h.service.ShortenURL(c.Request.Context(), req.URL, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := ShortenResponse{
		ShortURL:  code,
		ExpiresAt: time.Now().Add(5 * 365 * 24 * time.Hour),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("shortCode")
	url, err := h.service.ResolveURL(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Redirect(http.StatusFound, url)
}

func (h *Handler) GetStats(c *gin.Context) {
	code := c.Param("shortCode")
	url, err := h.service.GetURLStats(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	resp := StatsResponse{
		OriginalURL:  url.OriginalURL,
		ShortURL:     url.ShortCode,
		AccessCount:  url.AccessCount,
		CreatedAt:    url.CreatedAt,
		LastAccessed: url.LastAccessed,
	}
	c.JSON(http.StatusOK, resp)
}
