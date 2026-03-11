package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gianglt2198/pizza-app/internal/cache"
	"github.com/gianglt2198/pizza-app/internal/model"
	"github.com/gianglt2198/pizza-app/internal/repository"
)

type OrderHandler struct {
	orderRepo *repository.OrderRepository
	pizzaRepo *repository.PizzaRepository
	cache     *cache.RedisCache
}

func NewOrderHandler(
	orderRepo *repository.OrderRepository,
	pizzaRepo *repository.PizzaRepository,
	cache *cache.RedisCache,
) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo, pizzaRepo: pizzaRepo, cache: cache}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Calculate total price
	totalPrice, err := h.calculateTotal(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderRepo.Create(ctx, req, totalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.orderRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	orders, err := h.orderRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderRepo.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) calculateTotal(ctx context.Context, req model.CreateOrderRequest) (float64, error) {
	var total float64
	for _, item := range req.Items {
		pizza, err := h.pizzaRepo.FindByID(ctx, item.PizzaID)
		if err != nil {
			return 0, err
		}
		total += pizza.Price * float64(item.Quantity)
	}
	return total, nil
}
