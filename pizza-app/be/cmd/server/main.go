package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gianglt2198/pizza-app/config"
	"github.com/gianglt2198/pizza-app/internal/cache"
	"github.com/gianglt2198/pizza-app/internal/handlers"
	"github.com/gianglt2198/pizza-app/internal/repository"
)

func main() {
	// Load config
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}
	log.Println("✅ PostgreSQL connected")

	// Connect Redis
	rdb := cache.NewRedisCache(cfg.RedisAddr(), cfg.Redis.RedisPassword)
	if err = rdb.Ping(ctx); err != nil {
		log.Fatalf("failed to ping redis: %v", err)
	}
	log.Println("✅ Redis connected")

	// Repositories
	pizzaRepo := repository.NewPizzaRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)

	// Handlers
	healthHandler := handlers.NewHealthHandler(rdb)
	pizzaHandler := handlers.NewPizzaHandler(pizzaRepo, rdb)
	orderHandler := handlers.NewOrderHandler(orderRepo, pizzaRepo, rdb)

	// Router
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Routes
	r.GET("/health", healthHandler.Check)

	v1 := r.Group("/api/v1")
	{
		pizzas := v1.Group("/pizzas")
		{
			pizzas.GET("", pizzaHandler.List)
			pizzas.GET("/available", pizzaHandler.ListAvailable)
			pizzas.GET("/:id", pizzaHandler.GetByID)
		}

		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.Create)
			orders.GET("", orderHandler.List)
			orders.GET("/:id", orderHandler.GetByID)
			orders.PATCH("/:id/status", orderHandler.UpdateStatus)
		}
	}

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.App.ServerPort,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("🍕 Pizza Hub running on :%s", cfg.App.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
