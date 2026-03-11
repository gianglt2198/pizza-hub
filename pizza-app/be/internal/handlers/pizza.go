package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/gianglt2198/pizza-app/internal/cache"
	"github.com/gianglt2198/pizza-app/internal/model"
	"github.com/gianglt2198/pizza-app/internal/repository"
)

type PizzaHandler struct {
	repo  *repository.PizzaRepository
	cache *cache.RedisCache
}

func NewPizzaHandler(repo *repository.PizzaRepository, cache *cache.RedisCache) *PizzaHandler {
	return &PizzaHandler{repo: repo, cache: cache}
}

func (h *PizzaHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	cacheKey := "pizzas:all"

	var pizzas []model.Pizza
	if err := h.cache.Get(ctx, cacheKey, &pizzas); err == nil {
		c.Header("X-Cache", "HIT")
		c.JSON(http.StatusOK, pizzas)
		return
	}

	pizzas, err := h.repo.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.cache.Set(ctx, cacheKey, pizzas, 5*time.Minute)
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, pizzas)
}

func (h *PizzaHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx := c.Request.Context()
	cacheKey := "pizza:" + strconv.Itoa(id)

	var pizza model.Pizza
	if err := h.cache.Get(ctx, cacheKey, &pizza); err == nil {
		c.Header("X-Cache", "HIT")
		c.JSON(http.StatusOK, pizza)
		return
	}

	p, err := h.repo.FindByID(ctx, id)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "pizza not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.cache.Set(ctx, cacheKey, p, 5*time.Minute)
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, p)
}

func (h *PizzaHandler) ListAvailable(c *gin.Context) {
	ctx := c.Request.Context()
	cacheKey := "pizzas:available"

	var pizzas []model.Pizza
	if err := h.cache.Get(ctx, cacheKey, &pizzas); err == nil {
		c.Header("X-Cache", "HIT")
		c.JSON(http.StatusOK, pizzas)
		return
	}

	pizzas, err := h.repo.FindAvailable(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.cache.Set(ctx, cacheKey, pizzas, 5*time.Minute)
	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, pizzas)
}
