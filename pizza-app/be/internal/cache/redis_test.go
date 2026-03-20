package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_SetAndGet(t *testing.T) {
	s := miniredis.RunT(t)

	cache := NewRedisCache(s.Addr(), "")
	ctx := context.Background()

	err := cache.Set(ctx, "test_key", "test_value", 1*time.Minute)
	assert.NoError(t, err)

	var val string
	err = cache.Get(ctx, "test_key", &val)
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)
}

func TestRedisCache_Delete(t *testing.T) {
	s := miniredis.RunT(t)

	cache := NewRedisCache(s.Addr(), "")
	ctx := context.Background()

	cache.Set(ctx, "test_key", "test_value", 1*time.Minute)

	err := cache.Delete(ctx, "test_key")
	assert.NoError(t, err)

	var val string
	err = cache.Get(ctx, "test_key", &val)
	assert.Error(t, err) // Should error because not found
}
