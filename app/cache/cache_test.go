package cache

import (
	"example/web-service-gin/testUtils"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

const (
	serviceName = "testService"
)

func setupTestRedis(t *testing.T) (*miniredis.Miniredis, Cacher) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}

	port, err := strconv.Atoi(mr.Port())
	if err != nil {
		t.Fatalf("Failed to convert miniredis port to int: %v", err)
	}

	cfg := RedisClientConfig{
		Host: mr.Host(),
		Port: port,
	}

	return mr, NewCacher(cfg)
}

func TestGet(t *testing.T) {
	mr, cacher := setupTestRedis(t)
	defer mr.Close()

	key := "testKey"
	value := "testValue"

	// Set a value in Redis
	mr.Set(key, value)

	t.Run("Test Get", func(t *testing.T) {
		ctx := testUtils.CreateTestContext()
		result, err := cacher.Get(serviceName, ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Test Get with non-existent key", func(t *testing.T) {
		ctx := testUtils.CreateTestContext()
		result, err := cacher.Get(serviceName, ctx, "nonExistentKey")
		assert.Equal(t, ErrCacheMiss, err)
		assert.Equal(t, "", result)
	})

	t.Run("Test Get with Redis error", func(t *testing.T) {
		ctx := testUtils.CreateTestContext()
		mr.Close()
		_, err := cacher.Get(serviceName, ctx, key)
		assert.Equal(t, ErrCacheGeneric, err)
	})

}

func TestSet(t *testing.T) {
	mr, cacher := setupTestRedis(t)
	defer mr.Close()

	ctx := testUtils.CreateTestContext()
	key := "testKey"
	value := "testValue"
	expiration := time.Minute

	// Test Set
	err := cacher.Set(serviceName, ctx, key, value, expiration)
	assert.NoError(t, err)

	// Verify the value was set in Redis
	result, err := mr.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Verify expiration was set (approximately)
	ttl := mr.TTL(key)
	assert.InDelta(t, expiration.Seconds(), ttl.Seconds(), 1)
}
