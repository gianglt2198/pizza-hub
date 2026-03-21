//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPizzaRepository_FindAll(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()

	requireSeedExists(t, pool)

	repo := NewPizzaRepository(pool)
	pizzas, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, pizzas)
}

func TestPizzaRepository_FindByID(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()

	id := insertTestPizza(t, pool, "FindByID Pizza", 9.5, true)

	repo := NewPizzaRepository(pool)
	p, err := repo.FindByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, p.ID)
	require.Equal(t, "FindByID Pizza", p.Name)
}

func TestPizzaRepository_FindAvailable(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()

	// Ensure we have at least 1 unavailable, 1 available
	_ = insertTestPizza(t, pool, "Avail TRUE", 10, true)
	_ = insertTestPizza(t, pool, "Avail FALSE", 10, false)

	repo := NewPizzaRepository(pool)
	pizzas, err := repo.FindAvailable(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, pizzas)

	for _, p := range pizzas {
		require.True(t, p.Available, "FindAvailable returned unavailable pizza id=%d", p.ID)
	}
}