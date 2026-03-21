//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// In CI we will set TEST_DATABASE_URL
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// Fallback for local runs (matching docker-compose ports)
		dsn = "postgres://pizza:pizza123@localhost:5432/pizzahub?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create pg pool: %v", err)
	}

	// Wait until DB is ready (important for CI)
	deadline := time.Now().Add(20 * time.Second)
	for {
		if err := pool.Ping(ctx); err == nil {
			break
		}
		if time.Now().After(deadline) {
			pool.Close()
			t.Fatalf("postgres not ready before deadline")
		}
		time.Sleep(500 * time.Millisecond)
	}

	return pool
}

func truncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		TRUNCATE TABLE order_items, orders RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("failed to truncate: %v", err)
	}
}

func insertTestPizza(t *testing.T, pool *pgxpool.Pool, name string, price float64, available bool) int {
	t.Helper()

	var id int
	err := pool.QueryRow(context.Background(), `
		INSERT INTO pizzas (name, description, price, category, available)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, name, "test", price, "classic", available).Scan(&id)
	if err != nil {
		t.Fatalf("insert pizza failed: %v", err)
	}
	return id
}

func requireSeedExists(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	var cnt int
	if err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM pizzas`).Scan(&cnt); err != nil {
		t.Fatalf("count pizzas failed: %v", err)
	}
	if cnt == 0 {
		t.Fatalf("expected seeded pizzas, got 0")
	}
	fmt.Println("seed pizzas:", cnt)
}