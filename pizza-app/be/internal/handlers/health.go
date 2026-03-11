package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gianglt2198/pizza-app/internal/cache"
)

type HealthHandler struct {
	cache *cache.RedisCache
}

func NewHealthHandler(cache *cache.RedisCache) *HealthHandler {
	return &HealthHandler{cache: cache}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx := c.Request.Context()
	status := gin.H{
		"status":  "ok",
		"service": "pizza-hub",
	}

	if err := h.cache.Ping(ctx); err != nil {
		status["redis"] = "unhealthy: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	status["redis"] = "healthy"
	c.JSON(http.StatusOK, status)
}
