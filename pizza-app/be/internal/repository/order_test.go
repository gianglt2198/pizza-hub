//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gianglt2198/pizza-app/internal/model"
)

func TestOrderRepository_Create_And_FindByID(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	truncateTables(t, pool)

	// ensure a pizza exists to reference
	pizzaID := insertTestPizza(t, pool, "Order Pizza", 12.34, true)

	repo := NewOrderRepository(pool)

	req := model.CreateOrderRequest{
		CustomerName: "Alice",
		Items: []struct {
			PizzaID  int `json:"pizza_id" binding:"required"`
			Quantity int `json:"quantity" binding:"required,min=1"`
		}{
			{PizzaID: pizzaID, Quantity: 2},
		},
	}

	// totalPrice is computed in handler; repository receives it already
	order, err := repo.Create(context.Background(), req, 24.68)
	require.NoError(t, err)
	require.NotZero(t, order.ID)
	require.Len(t, order.Items, 1)
	require.Equal(t, model.StatusPending, order.Status)

	loaded, err := repo.FindByID(context.Background(), order.ID)
	require.NoError(t, err)
	require.Equal(t, order.ID, loaded.ID)
	require.Len(t, loaded.Items, 1)
	require.Equal(t, pizzaID, loaded.Items[0].PizzaID)
	require.Equal(t, 2, loaded.Items[0].Quantity)
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	truncateTables(t, pool)

	pizzaID := insertTestPizza(t, pool, "Status Pizza", 10, true)

	repo := NewOrderRepository(pool)
	req := model.CreateOrderRequest{
		CustomerName: "Bob",
		Items: []struct {
			PizzaID  int `json:"pizza_id" binding:"required"`
			Quantity int `json:"quantity" binding:"required,min=1"`
		}{
			{PizzaID: pizzaID, Quantity: 1},
		},
	}

	order, err := repo.Create(context.Background(), req, 10)
	require.NoError(t, err)

	updated, err := repo.UpdateStatus(context.Background(), order.ID, model.StatusPreparing)
	require.NoError(t, err)
	require.Equal(t, model.StatusPreparing, updated.Status)
	require.True(t, updated.UpdatedAt.After(updated.CreatedAt) || updated.UpdatedAt.Equal(updated.CreatedAt))
}

func TestOrderRepository_FindAll(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	truncateTables(t, pool)

	pizzaID := insertTestPizza(t, pool, "List Pizza", 10, true)
	repo := NewOrderRepository(pool)

	for i := 0; i < 2; i++ {
		req := model.CreateOrderRequest{
			CustomerName: "User",
			Items: []struct {
				PizzaID  int `json:"pizza_id" binding:"required"`
				Quantity int `json:"quantity" binding:"required,min=1"`
			}{
				{PizzaID: pizzaID, Quantity: 1},
			},
		}
		_, err := repo.Create(context.Background(), req, 10)
		require.NoError(t, err)
	}

	orders, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	require.Len(t, orders, 2)
}