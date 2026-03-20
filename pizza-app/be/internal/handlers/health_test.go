package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockCache struct {
	PingFunc   func(ctx context.Context) error
	GetFunc    func(ctx context.Context, key string, dest any) error
	SetFunc    func(ctx context.Context, key string, value any, ttl time.Duration) error
	DeleteFunc func(ctx context.Context, keys ...string) error
}

func (m *MockCache) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}
func (m *MockCache) Get(ctx context.Context, key string, dest any) error {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key, dest)
	}
	return nil
}
func (m *MockCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, ttl)
	}
	return nil
}
func (m *MockCache) Delete(ctx context.Context, keys ...string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, keys...)
	}
	return nil
}

func TestHealthCheck_Healthy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCache := &MockCache{
		PingFunc: func(ctx context.Context) error {
			return nil // Healthy
		},
	}

	handler := NewHealthHandler(mockCache)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	handler.Check(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"redis":"healthy"`)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestHealthCheck_UnhealthyCache(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCache := &MockCache{
		PingFunc: func(ctx context.Context) error {
			return errors.New("redis connection refused")
		},
	}

	handler := NewHealthHandler(mockCache)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	handler.Check(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"redis":"unhealthy: redis connection refused"`)
}
