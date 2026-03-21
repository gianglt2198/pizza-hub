package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func TestHealthCheck_Healthy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockCache := new(MockCache)
	mockCache.On("Ping", context.Background()).Return(nil)
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

	mockCache := new(MockCache)
	mockCache.On("Ping", context.Background()).Return(errors.New("redis connection refused"))
	handler := NewHealthHandler(mockCache)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	handler.Check(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"redis":"unhealthy: redis connection refused"`)
}
