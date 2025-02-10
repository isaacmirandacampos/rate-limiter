package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/isaacmirandacampos/rate-limiter/internal/core"
	"github.com/isaacmirandacampos/rate-limiter/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// Helper function to start a Redis container and return a connection pool
func setupRedisTestContainer(t *testing.T) (*redigo.Pool, func()) {
	// Start a Redis test container
	ctx := context.Background()
	redisContainer, err := redis.Run(ctx,
		"redis:7.4.2-alpine",
	)
	assert.NoError(t, err)

	// Get Redis connection info
	redisHost, err := redisContainer.Host(ctx)
	assert.NoError(t, err)
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	assert.NoError(t, err)
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())

	redisPool := &redigo.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", redisAddr)
		},
	}

	// Cleanup function to terminate Redis container after tests
	cleanup := func() {
		redisContainer.Terminate(ctx)
	}

	return redisPool, cleanup
}

// Helper function to test rate limiting
func testRateLimiter(t *testing.T, rateLimiter *middleware.RateLimiter, headers map[string]string, expectedPasses int) {
	// Mock handler function
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimiter.RateLimiterMiddleware(mockHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	for i := 1; i <= expectedPasses; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should pass", i)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request %d should be rate limited", expectedPasses+1)
}

// Test rate limiter based on IP
func TestRateLimiterMiddlewareByIp(t *testing.T) {
	redisPool, cleanup := setupRedisTestContainer(t)
	defer cleanup()

	// Initialize rate limiter
	repo := core.NewRedisRateLimiterRepository(redisPool)
	rateLimiter := middleware.NewRateLimiter(core.NewRateLimiterByIp(repo, 3), 2, 1)

	// Run test with IP-based rate limiting
	testRateLimiter(t, rateLimiter, nil, 2)
}

// Test rate limiter based on API Key
func TestRateLimiterMiddlewareByApiKey(t *testing.T) {
	redisPool, cleanup := setupRedisTestContainer(t)
	defer cleanup()

	// Initialize rate limiter
	repo := core.NewRedisRateLimiterRepository(redisPool)
	rateLimiter := middleware.NewRateLimiter(core.NewRateLimiterByIp(repo, 3), 1, 2)

	// Run test with API key-based rate limiting
	testRateLimiter(t, rateLimiter, map[string]string{"API_KEY": "1234"}, 2)
}
